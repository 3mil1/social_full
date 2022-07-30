package service

import (
	"errors"
	"social-network/internal/dto"
	"social-network/internal/entity"
)

type GroupService struct {
	repo     GroupRepo
	follower FollowerRepo
}

func NewGroupService(r GroupRepo, follower FollowerRepo) *GroupService {
	return &GroupService{
		r,
		follower,
	}
}

func (service *GroupService) AddNewGroup(group entity.Group) (int, error) {
	return service.repo.AddNewGroup(group)
}

func (service *GroupService) GetAllGroups() ([]dto.GroupShortInfo, error) {
	return service.repo.GetAllGroups()
}

func (service *GroupService) GetAllMyCreatedGroups(loggedInUser string) ([]dto.GroupShortInfo, error) {
	return service.repo.GetAllMyCreatedGroups(loggedInUser)
}

func (service *GroupService) GetAllMyJoinedGroups(loggedInUser string) ([]dto.GroupShortInfo, error) {
	return service.repo.GetAllMyJoinedGroups(loggedInUser)
}

func (service *GroupService) GetOneGroupInfo(loggedInUser string, groupId int) (dto.GroupDetailedInfo, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUser, groupId)
	if err != nil {
		return dto.GroupDetailedInfo{}, err
	}
	if isMember {
		return service.repo.GetOneGroupInfo(groupId)
	}
	groupInfo, err := service.repo.GetOneGroupInfo(groupId)
	if err != nil {
		return dto.GroupDetailedInfo{}, err
	}
	return groupInfo, errors.New("user is not a group member")
}

func (service *GroupService) AddNewGroupPost(post entity.GroupPost) error {
	isMember, err := service.repo.CheckUserIsGroupMember(post.CreatorId, post.GroupId)
	if err != nil {
		return err
	}
	if isMember {
		return service.repo.AddNewGroupPost(post)
	}
	return errors.New("user is not a group member")
}

func (service *GroupService) GetOneGroupAllPosts(loggedInUser string, groupId int) ([]dto.GroupPostReply, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUser, groupId)
	if err != nil {
		return nil, err
	}
	if isMember {
		return service.repo.GetOneGroupAllPosts(groupId)
	}
	return nil, errors.New("user is not a group member")
}

func (service *GroupService) GetOnePostAndComments(loggedInUserId string, groupId, postId int) (dto.GroupOnePostAndComments, error) {
	var result dto.GroupOnePostAndComments
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUserId, groupId)
	if err != nil {
		return result, err
	}
	if isMember {
		post, err := service.repo.GetOnePost(groupId, postId)
		if err != nil {
			return result, err
		}
		comments, err := service.repo.GetPostComments(groupId, postId)
		if err != nil {
			return result, err
		}
		result.Post = post
		result.Comments = comments

		return result, nil
	}
	return result, errors.New("user is not a group member")
}

func (service *GroupService) InviteUserToGroup(loggedInUserId string, invitation dto.GroupInvitation) (dto.UserStatusInGroup, error) {
	var reply dto.UserStatusInGroup
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUserId, invitation.GroupId)
	if err != nil {
		return reply, err
	}
	if isMember {
		//new request in not in group already
		status, err := service.repo.CheckUserStatusInGroup(invitation.TargetId, invitation.GroupId)
		if err != nil {
			return reply, err
		}
		reply.Status = status
		if status < 0 {
			// 0 -> pending -> admin or group member invite user and USER has to reply
			// 1 -> Accept
			// 2 -> Decline
			// 3->  request -> user makes join request and group ADMIN needs to reply to it
			err := service.repo.AddGroupMember(invitation.GroupId, invitation.TargetId, "member", 0)
			if err != nil {
				return reply, err
			}
			reply.Status = 0
		}
		return reply, nil
	}
	return reply, errors.New("loggedInUser is not a group member")
}

func (service *GroupService) ReplyToGroupInvitation(loggedInUserId string, reply dto.GroupInvitationReply) error {
	if reply.Status != 2 {
		return service.repo.ReplyToGroupInvitation(loggedInUserId, reply.Status, reply.GroupId)
	}
	return service.repo.DeleteInvitation(loggedInUserId, reply.GroupId)
}

func (service *GroupService) RequestGroupAccess(loggedInUserId string, groupId int) (bool, dto.UserStatusInGroup, error) {
	var reply dto.UserStatusInGroup
	newRequest := false
	status, err := service.repo.CheckUserStatusInGroup(loggedInUserId, groupId)
	if err != nil {
		return newRequest, reply, err
	}
	reply.Status = status
	
	if status < 0 {
		newRequest = true
		// 0 -> pending -> admin or group member invite user and USER has to reply
		// 1 -> Accept
		// 2 -> Decline
		// 3->  request -> user makes join request and group ADMIN needs to reply to it
		err = service.repo.AddGroupMember(groupId, loggedInUserId, "member", 3)
		if err != nil {
			return newRequest, reply, err
		}
		reply.Status = 3
		return newRequest, reply, nil
	} else if status == 0{
		err =  service.repo.UpdateGroupMember(groupId, loggedInUserId, "member", 1)
		if err != nil{
			return newRequest, reply, err
		}
		reply.Status = 1
	}
	return newRequest, reply, nil
}

func (service *GroupService) ReplyToGroupAccessRequest(loggedInUserId string, reply dto.GroupAccessRequestReply) error {
	admin, err := service.repo.CheckUserRoleInGroup(loggedInUserId, reply.GroupId)
	if err != nil {
		return err
	}
	if !admin {
		return errors.New("loggedInUser is not an admin in group")
	}
	//if loggedInUser is admin
	if reply.Status == 2 {
		return service.repo.DeleteInvitation(reply.TargetId, reply.GroupId)
	}
	return service.repo.ReplyToGroupInvitation(reply.TargetId, reply.Status, reply.GroupId)
}

func (service *GroupService) CreateEvent(event entity.GroupEventEntity) (int, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(event.UserId, event.GroupId)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, errors.New("user is not a group member")
	}
	return service.repo.CreateEvent(event)
}

func (service *GroupService) AddEventParticipant(userId string, event dto.EventParticipant) error {
	groupId, err := service.repo.GetGroupIdFromEvent(event.EventId)
	if err != nil {
		return err
	}
	if groupId > 0 {
		isMember, err := service.repo.CheckUserIsGroupMember(userId, groupId)
		if err != nil {
			return err
		}
		if !isMember {
			return errors.New("user is not a group member")
		}
		//check that user does not have an added participation already (option is -1)
		option, err := service.repo.CheckUserGoingInEvent(event.EventId, userId)
		if err != nil {
			return err
		}
		if option < 0 {
			// no such listing-> add new
			return service.repo.AddEventParticipant(userId, event.EventId, event.Option)
		} else {
			return service.repo.UpdateEventParticipant(userId, event.EventId, event.Option)
		}
	}
	return errors.New("no such event or group presented")
}

func (service *GroupService) GetAllGroupEvents(loggedInUserId string, groupId int) ([]dto.GroupEventReply, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUserId, groupId)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a group member")
	}
	return service.repo.GetAllGroupEvents(groupId)
}

func (service *GroupService) GetMyCreatedEvents(loggedInUserId string) ([]dto.GroupEventReply, error) {
	return service.repo.GetMyCreatedEvents(loggedInUserId)
}

func (service *GroupService) GetMyJoinedEvents(loggedInUserId string) ([]dto.GroupEventReply, error) {
	return service.repo.GetMyJoinedEvents(loggedInUserId)
}

// func (service *GroupService) AddNewGroupMessage(message entity.GroupMessage) (entity.GroupMessage, error) {
// 	isMember, err := service.repo.CheckUserIsGroupMember(message.UserId, message.GroupId)
// 	if err != nil {
// 		return entity.GroupMessage{}, err
// 	}
// 	if !isMember {
// 		return entity.GroupMessage{}, errors.New("user is not a group member")
// 	}
// 	return service.repo.AddNewGroupMessage(message)
// }

func (service *GroupService) GetOneGroupAllMessages(loggedInUserId string, groupId, skip, limit int) ([]dto.GroupMessageReply, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUserId, groupId)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a group member")
	}
	return service.repo.GetOneGroupAllMessages(loggedInUserId, groupId, skip, limit)
}

func (service *GroupService) GetGroupAdmin(groupId int) (string, error) {
	return service.repo.GetGroupAdmin(groupId)
}

func (service *GroupService) GetPendingJoinRequests(loggedInUser string, groupId int) ([]dto.GroupAccessRequestUser, error) {
	admin, err := service.repo.CheckUserRoleInGroup(loggedInUser, groupId)
	if err != nil {
		return nil, err
	}
	if !admin {
		return nil, errors.New("loggedInUser is not an admin in group")
	}
	//if loggedInUser is admin
	return service.repo.GetPendingJoinRequests(groupId)
}

func (service *GroupService) GetFriendsNotInGroup(loggedInUserId string, groupId int) ([]dto.FollowerUserSmall, error) {
	isMember, err := service.repo.CheckUserIsGroupMember(loggedInUserId, groupId)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a group member")
	}
	var result []dto.FollowerUserSmall
	myFriends, err := service.follower.GetAllUsersIFollow(loggedInUserId)
	if err != nil {
		return nil, err
	}

	for _, oneFriend := range myFriends {
		requestExist, err := service.repo.CheckUserStatusInGroup(oneFriend.UserId, groupId)
		if err != nil {
			return nil, err
		}
		if requestExist < 0 {
			result = append(result, oneFriend)
		}
	}
	return result, nil
}

func (service *GroupService) GetAllGroupMembersExceptMe(loggedInUser string, groupId int) ([]string, error) {
	return service.repo.GetAllGroupMembersExceptMe(loggedInUser, groupId)
}
