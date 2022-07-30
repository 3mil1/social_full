package dto

import "social-network/internal/entity"

type GroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GroupRequestToEntity(loggedInUser string, group GroupRequest) entity.Group {
	return entity.Group{
		Title:       group.Title,
		Description: group.Description,
		Creator:     loggedInUser,
	}
}

type GroupDetailedInfo struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	CreatorId        string `json:"creator_id"`
	CreatorFirstName string `json:"creator_first_name"`
	CreatorLastName  string `json:"creator_last_name"`
	Members          []PrivateProfileResponse
}

type GroupShortInfo struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Members          int	`json:"members"`
}


type GroupPost struct {
	GroupId  int    `json:"group_id"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Image    string `json:"image"`
	ParentId int    `json:"parent_id"`
}

func GroupPostToEntity(loggedInUser string, post GroupPost) entity.GroupPost {
	return entity.GroupPost{
		GroupId:   post.GroupId,
		Title:     post.Subject,
		CreatorId: loggedInUser,
		Content:   post.Content,
		Image:     post.Image,
		ParentId:  post.ParentId,
	}
}

type GroupPostReply struct {
	PostId        int    `json:"post_id"`
	UserId        string `json:"user_id"`
	UserFirstName string `json:"user_firstname"`
	UserLastName  string `json:"User_lastname"`
	Subject       string `json:"subject"`
	Content       string `json:"content"`
	Image         string `json:"image"`
	ParentId      int    `json:"parent_id"`
	CreatedAt     string `json:"created_at"`
}
type GroupOnePostAndComments struct{
	Post  		GroupPostReply
	Comments 	[]GroupPostReply
}

type GroupInvitation struct{
	GroupId 	int		`json:"group_id"`
	TargetId 	string 	`json:"target_id"`
}

type GroupInvitationReply struct{
	GroupId 	int		`json:"group_id"`
	ActorId		string	`json:"actor_id"`
	Status 		int 	`json:"status"`
}

type GroupId struct{
	GroupId 	int		`json:"group_id"`
}

type UserStatusInGroup struct{
	Status 		int 	`json:"status"`
}

type GroupEvent struct{
	GroupId 		int		`json:"group_id"`
	Title 			string 	`json:"title"`
	Description 	string 	`json:"description"`
	Day				string 	`json:"day"`
	Time 			string	`json:"time"`
	Going 			int		`json:"going_status"`
}

func GroupEventToEntity(loggedInUser string, event GroupEvent) entity.GroupEventEntity{
	return entity.GroupEventEntity{
		GroupId: event.GroupId,
		UserId: loggedInUser,
		Title: event.Title,
		Description: event.Description,
		EventDate: event.Day +" "+ event.Time,
		GoingStatus: event.Going,
	}
}

type EventParticipant struct{
	EventId		int `json:"event_id"`
	Option 		int `json:"option"`
}

type GroupAccessRequestReply struct{
	GroupId 	int		`json:"group_id"`
	TargetId 	string 	`json:"target_id"`
	Status 		int 	`json:"status"`
}

type GroupEventReply struct{
	Id 				int 	`json:"event_id"`
	GroupId 		int		`json:"group_id"`
	CreatorId 		string 	`json:"creator_id"`
	CreatorFirstName string	`json:"creator_firstname"`
	CreatorLastName string	`json:"creator_lastname"`
	Title 			string 	`json:"title"`
	Description 	string 	`json:"description"`
	Day				string 	`json:"day"`
	Time 			string	`json:"time"`
	CreatedAt 		string 	`json:"created_at"`
}

type GroupMessageReply struct{
	MessageId 		int 	`json:"message_id,omitempty"`
	GroupId 		int		`json:"group_id,omitempty"`
	FromId 			string 	`json:"from,omitempty"`
	FirstName 		string	`json:"first_name,omitempty"`
	LastName 		string	`json:"last_name,omitempty"`
	Content 		string 	`json:"content,omitempty"`
	CreatedAt 		string 	`json:"created_at,omitempty"`
	Seen 			bool	`json:"seen,omitempty"`
}

type GroupAccessRequestUser struct{
	GroupId 		int		`json:"group_id"`
	UserId 			string 	`json:"user_id"`
	UserFirstName	string 	`json:"user_firstname"`
	UserLastName	string 	`json:"user_lastname"`
	Role 			string  `json:"role"`
	Status 			int 	`json:"status"`
}