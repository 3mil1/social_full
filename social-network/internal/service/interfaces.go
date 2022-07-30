package service

import (
	"social-network/internal/dto"
	"social-network/internal/entity"
)

type (
	User interface {
		AddUser(data dto.UserRequestBody) (*dto.UserResponse, error)
		SignIn(data dto.SignInRequestBody, ip, userAgent string) (*dto.TokenResponse, error)
		RefreshToken(refreshTokenFromRequest dto.RefreshTokenRequestBody, ip, userAgent string) (*dto.TokenResponse, error)
		GetUserByID(userID string) (*dto.UserResponse, error)
		SignOut(userID, ip, userAgent string) error
		UpdateUser(user dto.UserUpdate, loggedInUser string) error
		GetMyFollowerProfile(loggedInUserId, frinedId string) (*dto.UserResponse, error)
		GetAllUsers() ([]dto.PrivateProfileResponse, error)
	}

	UserRepo interface {
		AddUser(model *entity.User) (*entity.User, error)
		GetUserByID(id string) (*entity.User, error)
		GetUserByEmail(email string) (*entity.User, error)
		AddRefreshToken(refreshToken entity.RefreshTokenDB) error
		DeleteRefreshToken(refreshToken entity.RefreshTokenDB) (*entity.RefreshTokenDB, error)
		DeleteSession(userID string, ip, userAgent string) error
		UpdateUser(user dto.UserUpdate, loggedInUser string) error
		GetUserStatusByID(requestedUserId string) (bool, error)
		Get2UsersConnectionStatus(loggedInUserId, requestedUserId string) (int, error)
		GetAllUsers() ([]dto.PrivateProfileResponse, error)
	}

	FollowerServe interface {
		GetAllUsersIFollow(loggedInUserid, someUserId string) ([]dto.FollowerUserSmall, error)
		AddNewFollower(follower entity.Follower) (int, error)
		UpdateFollower(follower entity.Follower) error
		GetAllUsersFollowsMe(targetId string) ([]dto.FollowerUserSmall, error)
		GetUsersFollowsMeAcceptedStatusOnly(loggedInUserId, targetId string) ([]dto.FollowerUserSmall, error)
		UpdateFollowRequest(follower entity.Follower) error
		DeleteFollower(loggedInUserId, userId string) error
		GetChatList(loggedInUserId string) ([]dto.FollowerUserSmall, error)
	}

	FollowerRepo interface {
		GetAllUsersIFollow(someUserId string) ([]dto.FollowerUserSmall, error)
		CheckFollowRequest(someUserId, targetId string) (int, error)
		AddNewFollower(follower entity.Follower) error
		UpdateFollower(follower entity.Follower) error
		GetUsersFollowsMeAcceptedStatusOnly(targetId string) ([]dto.FollowerUserSmall, error)
		GetAllUsersFollowsMe(targetId string) ([]dto.FollowerUserSmall, error)
		DeleteFollower(loggedInUserId, userId string) error
	}

	PostServe interface {
		AddNewPost(post dto.PostReceive) (dto.PostReply, error)
		GetAllUserPosts(loggedInUserId, requestedUserId string) ([]dto.PostReply, error)
		GetOnePostWithComments(loggedInUserId string, postId int) (dto.PostAndComments, error)
		PostWithComments(postId int) (dto.PostAndComments, error)
		GetPostOwner(postId int) (string, error)
		GetAllPosts(loggedInUserId string) ([]dto.PostReply, error)
	}

	PostRepo interface {
		AddNewPost(post entity.UserPost) (int, error)
		AddNewPostAccess(postId int, oneUserId string) error
		GetAllUserPosts(loggedInUserId, requestedUserId string) ([]dto.PostReply, error)
		GetOnePostsComments(postId int) ([]dto.PostReply, error)
		GetPostStatusByPostId(postId int) (int, string, error)
		GetPostAuthorByPostId(postId int) (string, error)
		CheckPostAccessByPostIdAndUserId(postId int, loggedInUser string) (int, error)
		GetOnePostByPostId(postId int) (dto.PostReply, error)
		GetPostOwner(postId int) (string, error)
		GetAllPosts(loggedInUserId string) ([]dto.PostReply, error)
	}

	ChatServe interface {
		GetMessages(sender, receiver string, skip, limit int) ([]dto.ChatMessage, error)
		AddMessage(sender, receiver, messageContent string) ([]dto.WsResponse, error)
		AddGroupMessage(sender, groupId, message string) ([]dto.WsResponse, []dto.PrivateProfileResponse, error)
		MarkMessageAsSeen(loggedInUserId, target_id string)
	}

	ChatRepo interface {
		//CountMessages(sender, receiver string) (int, error)
		GetMessages(sender, receiver string, skip, limit int) ([]dto.ChatMessage, error)
		AddMessage(chatMessage *entity.ChatMessage, sender, messageContent string, chatID int, receiver string) error
		AddChat(sender, receiver, messageContent string) (*entity.ChatMessage, error)
		AddGroupMessage(sender, groupId, messageContent string) (*entity.GroupChatMessage, error)
		AddGroupChatReceiver(messageId int, receiverId string, read bool) error
		MarkPersonalMessageAsSeen(sender, receiver string)
		MarkGroupMessagesAsSeen(sender, receiver string)
	}

	GroupServe interface {
		AddNewGroup(groupEntity entity.Group) (int, error)
		GetAllGroups() ([]dto.GroupShortInfo, error)
		GetAllMyCreatedGroups(loggedInUser string) ([]dto.GroupShortInfo, error)
		GetAllMyJoinedGroups(loggedInUser string) ([]dto.GroupShortInfo, error)
		GetOneGroupInfo(loggedInUser string, groupId int) (dto.GroupDetailedInfo, error)
		AddNewGroupPost(postEntity entity.GroupPost) error
		GetOneGroupAllPosts(loggedInUserId string, groupId int) ([]dto.GroupPostReply, error)
		GetOnePostAndComments(loggedInUserId string, groupId, postId int) (dto.GroupOnePostAndComments, error)
		InviteUserToGroup(loggedInUserId string, invitation dto.GroupInvitation) (dto.UserStatusInGroup, error)
		ReplyToGroupInvitation(loggedInUserId string, reply dto.GroupInvitationReply) error
		RequestGroupAccess(loggedInUserId string, groupId int) (bool, dto.UserStatusInGroup, error)
		ReplyToGroupAccessRequest(loggedInUserId string, reply dto.GroupAccessRequestReply) error
		CreateEvent(event entity.GroupEventEntity) (int,error)
		AddEventParticipant(userId string, event dto.EventParticipant) error
		GetAllGroupEvents(loggedInUserId string, groupId int) ([]dto.GroupEventReply, error)
		GetMyCreatedEvents(loggedInUserId string) ([]dto.GroupEventReply, error)
		GetMyJoinedEvents(loggedInUserId string) ([]dto.GroupEventReply, error)
		//AddNewGroupMessage(message entity.GroupMessage) (entity.GroupMessage, error)
		GetOneGroupAllMessages(sender string, groupId, skip, limit int) ([]dto.GroupMessageReply, error)
		GetGroupAdmin(groupId int) (string, error)
		GetPendingJoinRequests(loggedInUserId string, groupId int) ([]dto.GroupAccessRequestUser, error)
		GetFriendsNotInGroup(loggedInUserId string, groupId int) ([]dto.FollowerUserSmall, error)
		GetAllGroupMembersExceptMe(loggedInUserId string, groupId int) ([]string, error)
	}

	GroupRepo interface {
		AddNewGroup(groupEntity entity.Group) (int, error)
		GetAllGroups() ([]dto.GroupShortInfo, error)
		GetAllMyCreatedGroups(loggedInUser string) ([]dto.GroupShortInfo, error)
		GetAllMyJoinedGroups(loggedInUser string) ([]dto.GroupShortInfo, error)
		AddGroupMember(groupId int, userId, role string, status int) error
		GetGroupMemberNumber(groupId int) (int, error)
		GetGroupMemberDetails(groupId int) ([]dto.PrivateProfileResponse, error)
		GetOneGroupInfo(groupId int) (dto.GroupDetailedInfo, error)
		AddNewGroupPost(postEntity entity.GroupPost) error
		GetOneGroupAllPosts(groupId int) ([]dto.GroupPostReply, error)
		CheckUserIsGroupMember(loggedInUserId string, groupId int) (bool, error)
		GetOnePost(groupId, postId int) (dto.GroupPostReply, error)
		GetPostComments(groupId, postId int) ([]dto.GroupPostReply, error)
		ReplyToGroupInvitation(loggedInUser string, status, groupId int) error
		CheckUserStatusInGroup(loggedInUserId string, groupId int) (int, error)
		DeleteInvitation(loggedInUserId string, groupId int) error
		CreateEvent(event entity.GroupEventEntity) (int, error)
		AddEventParticipant(userId string, eventId, option int) error
		GetGroupIdFromEvent(eventId int) (int, error)
		CheckUserRoleInGroup(loggedInUser string, groupId int) (bool, error)
		CheckUserGoingInEvent(eventId int, userId string) (int, error)
		GetAllGroupEvents(groupId int) ([]dto.GroupEventReply, error)
		GetMyCreatedEvents(loggedInUserId string) ([]dto.GroupEventReply, error)
		GetMyJoinedEvents(loggedInUserId string) ([]dto.GroupEventReply, error)
		GetOneGroupAllMessages(loggedInUser string, groupId, skip, limit int) ([]dto.GroupMessageReply, error)
		GetGroupAdmin(groupId int) (string, error)
		GetPendingJoinRequests(groupId int) ([]dto.GroupAccessRequestUser, error)
		GetAllGroupMembersExceptMe(loggedInUserId string, groupId int) ([]string, error)
		UpdateEventParticipant(userId string, eventId, option int) error
		UpdateGroupMember(groupId int, userId, role string, status int) error
	}

	NotificationServe interface {
		AddNotification(actorID string, receiverID []string, notificationType, objectId int) error
		GetAllNotifications(id string) ([]dto.WsResponse, error)
		DeleteNotification(actorId, receiverId string, notification_type, objectID int) error
		UpdateNotification(loggedInUserId, notify_id, status string) error
		DeleteGroupChatNotification(loggedInUserId, group_id string) error 
	}

	NotificationRepo interface {
		AddNotificationReceiver(n entity.Notification) (*entity.Notification, error)
		AddNotificationOBJ(n entity.Notification) (*int, error)
		GetAllNotifications(id string) ([]dto.WsResponse, error)
		DeleteNotification(receiverId string, notifyId int) error
		FindNotificationId(actorId, receiverId string, notification_type, objectID int) ([]int, error)
		UpdateNotification(loggedInUserId, notify_id, status string) error
		DeleteGroupChatNotification(loggedInUserId, group_id string) error 
		CheckNotifyObj(n entity.Notification) (*int, error)
	}
)
