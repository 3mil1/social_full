package dto

import "social-network/internal/entity"

type PostReceive struct {
	UserId   string   `json:"user_id"`
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	Image    string   `json:"image"`
	ParentId int      `json:"parent_id"`
	Privacy  int      `json:"privacy"`
	Access   []string `json:"access"`
}

func PostReceivedToEntity(post PostReceive) entity.UserPost {
	return entity.UserPost{
		UserId:   post.UserId,
		Subject:  post.Subject,
		Content:  post.Content,
		Image:    post.Image,
		ParentId: post.ParentId,
		Privacy:  post.Privacy,
		Access:   post.Access,
	}
}

type PostReply struct {
	Id            int    `json:"id"`
	UserId        string `json:"user_id"`
	UserFirstName string `json:"user_firstname"`
	UserLastName  string `json:"user_lastname"`
	Subject       string `json:"subject"`
	Content       string `json:"content"`
	Image         string `json:"image"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Privacy       int    `json:"privacy"`
}

func PostEntityToPostReply(post entity.UserPost) PostReply {
	return PostReply{
		Id:      post.Id,
		UserId:  post.UserId,
		Subject: post.Subject,
		Content: post.Content,
		Image:   post.Image,
		Privacy: post.Privacy,
	}
}

type PostAndComments struct {
	Post     PostReply
	Comments []PostReply
}
