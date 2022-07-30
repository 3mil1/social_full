package service

import (
	"errors"
	"social-network/internal/dto"
)

type PostService struct {
	repo PostRepo
	user UserRepo
}

func NewPostService(r PostRepo, user UserRepo) *PostService {
	return &PostService{
		repo: r, 
		user: user,
	}
}

func (service *PostService) AddNewPost(post dto.PostReceive) (dto.PostReply, error) {

	newPost := dto.PostReceivedToEntity(post)

	//add new post into post table
	postId, err := service.repo.AddNewPost(newPost)
	newPost.Id = postId

	var reply dto.PostReply
	if err != nil {
		return reply, err
	}

	/* POST privacy rules:
	1 - PUBLIC
	2 - PRIVATE -> all followers af post author can see this post
	3 - DEEPLY PRIVATE-> only some of the post creator followers can see this post
	*/
	if len(newPost.Access) != 0 && newPost.Privacy == 3 {
		//add this dependency to the post_access table
		for _, oneUserId := range newPost.Access {
			err := service.repo.AddNewPostAccess(postId, oneUserId)
			if err != nil {
				return reply, err
			}
		}
	}
	reply, err = service.repo.GetOnePostByPostId(newPost.Id)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (service *PostService) GetAllUserPosts(loggedInUserId, requestedUserId string) ([]dto.PostReply, error) {
	return service.repo.GetAllUserPosts(loggedInUserId, requestedUserId)
}

func (service *PostService) GetOnePostWithComments(loggedInUserId string, postId int) (dto.PostAndComments, error) {
	var result dto.PostAndComments
	// check, that post is public/private/strictly_private
	access, owner, err := service.repo.GetPostStatusByPostId(postId)
	if err != nil {
		return result, err
	}
	//PUBLIC POST
	if access == 1 || owner == loggedInUserId {
		return service.PostWithComments(postId)
	}
	//PRIVATE POST
	if access == 2 {
		// check if loggedInUser has access to that post
		user_id, err := service.repo.GetPostAuthorByPostId(postId)
		if err != nil {
			return result, err
		}
		connection, err := service.user.Get2UsersConnectionStatus(loggedInUserId, user_id)
		if err != nil {
			return result, err
		}
		//IF APPROVED CONNECTION
		if connection == 1 {
			return service.PostWithComments(postId)
		} else {
			return result, nil
		}
	}
	//DeeplyPrivate post
	if access == 3 {
		id, err := service.repo.CheckPostAccessByPostIdAndUserId(postId, loggedInUserId)
		if err != nil {
			return result, err
		}
		if id > 0 {
			return service.PostWithComments(postId)
		} else {
			return result, errors.New("user has no access to this post")
		}
	}

	return result, nil
}

func (service *PostService) PostWithComments(postId int) (dto.PostAndComments, error) {
	var result dto.PostAndComments
	var post dto.PostReply
	var comments []dto.PostReply
	post, _ = service.repo.GetOnePostByPostId(postId)
	comments, _ = service.repo.GetOnePostsComments(postId)
	result.Post = post
	result.Comments = comments
	return result, nil
}

func (service *PostService) GetPostOwner(postId int) (string, error) {
	return service.repo.GetPostOwner(postId)
}

func (service *PostService) GetAllPosts(loggedInUserId string) ([]dto.PostReply, error) {
	return service.repo.GetAllPosts(loggedInUserId)
}