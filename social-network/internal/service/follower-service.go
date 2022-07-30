package service

import (
	"errors"
	"social-network/internal/dto"
	"social-network/internal/entity"
)

type FollowerService struct {
	repo FollowerRepo
	user UserRepo
}

func NewFollowerService(r FollowerRepo, u UserRepo) *FollowerService {
	return &FollowerService{
		repo: r,
		user: u,
	}
}

func (service *FollowerService) GetAllUsersIFollow(loggedInUserId, someUserId string) ([]dto.FollowerUserSmall, error) {
	if loggedInUserId == someUserId {
		return service.repo.GetAllUsersIFollow(loggedInUserId)
	}
	/*check that requested user is NOT PRIVATE profile
	FALSE - PUBLIC
	TRUE - Private
	*/
	private, err := service.user.GetUserStatusByID(someUserId)
	if err != nil {
		return nil, err
	}
	if !private {
		//status is public-> accept follower by default
		return service.repo.GetAllUsersIFollow(someUserId)
	} else {
		//check approved connection of user we need to check
		connection, err := service.user.Get2UsersConnectionStatus(loggedInUserId, someUserId)
		if err != nil {
			return nil, err
		}
		//IF APPROVED CONNECTION
		if connection == 1 {
			return service.repo.GetAllUsersIFollow(someUserId)
		} else {
			return nil, nil
		}
	}
}

func (service *FollowerService) AddNewFollower(follower entity.Follower) (int, error) {

	status, err := service.repo.CheckFollowRequest(follower.SourceId, follower.TargetId)
	if err != nil {
		return 0, err
	}
	private, err := service.user.GetUserStatusByID(follower.TargetId)
	if err != nil {
		return 0, err
	}
	if status >= 0 { //if connection already is requested/accepted/declined
		if !private{
			follower.Status = 1
		} else{
			follower.Status = status
		}
		err := service.repo.UpdateFollower(follower)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	/*check that requested user is NOT PRIVATE profile
	FALSE - PUBLIC
	TRUE - Private
	*/
	//new_connection to private user = > 1
	//new_connection to public user => 2
	new_connection := 1
	if !private {
		//status is public-> accept follower by default
		follower.Status = 1
		new_connection = 2
	} else {
		follower.Status = 0
	}
	err = service.repo.AddNewFollower(follower)
	if err != nil {
		return 0, err
	}
	//new connection to private user = > 1
	//new connection to public user => 2
	return new_connection, nil
}

func (service *FollowerService) UpdateFollower(follower entity.Follower) error {
	return service.repo.UpdateFollower(follower)
}

func (service *FollowerService) GetAllUsersFollowsMe(targetId string) ([]dto.FollowerUserSmall, error) {
	return service.repo.GetAllUsersFollowsMe(targetId)
}

func (service *FollowerService) GetUsersFollowsMeAcceptedStatusOnly(loggedInUserId, targetId string) ([]dto.FollowerUserSmall, error) {
	//check if targetId id PUBLIC
	/*check that requested user is NOT PRIVATE profile
	TRUE - PRIVATE
	FALSE - PUBLIC
	*/
	private, err := service.user.GetUserStatusByID(targetId)
	if err != nil {
		return nil, err
	}
	if !private {
		//status is public-> accept follower by default
		return service.repo.GetUsersFollowsMeAcceptedStatusOnly(targetId)
	} else {
		//check LoggedInUserIsAAcceptedFollower of targetId
		connected, err := service.user.Get2UsersConnectionStatus(loggedInUserId, targetId)
		if err != nil {
			return nil, err
		}
		//if loggedInUser have right to see requested person information (1 == approved)
		if connected == 1 {
			return service.repo.GetUsersFollowsMeAcceptedStatusOnly(targetId)
		}
		return nil, errors.New("you are not connected")
	}
}

func (service *FollowerService) UpdateFollowRequest(follower entity.Follower) error {
	middle := follower.SourceId
	follower.SourceId = follower.TargetId
	follower.TargetId = middle
	return service.UpdateFollower(follower)
}

func (service *FollowerService) DeleteFollower(loggedInUserId, userId string) error{
	return service.repo.DeleteFollower(loggedInUserId, userId)
}

func (service *FollowerService) GetChatList(loggedInUserId string) ([]dto.FollowerUserSmall, error){
	var friendList []dto.FollowerUserSmall
	list, err := service.repo.GetAllUsersIFollow(loggedInUserId) 
	if err != nil{
		return nil, err
	}
	private, err := service.user.GetUserStatusByID(loggedInUserId)
	if err != nil {
		return nil, err
	}
	//If logged in user is private profile-> select only users he is following as well
	if private{
		for _, oneUser := range list{
			status, err := service.repo.CheckFollowRequest(oneUser.UserId, loggedInUserId)
			if err != nil{
				return nil, err
			}
			if status == 1{
				friendList = append(friendList, oneUser)
			}
		}
	} else{
		//if logged in user is public account-> return all
		return list, nil
	}
	return friendList, nil
}