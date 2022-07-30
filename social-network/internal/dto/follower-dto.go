package dto

import "social-network/internal/entity"

type FollowerUserSmall struct{
	UserId 			string 	`json:"user_id,omitempty"`
	FirstName		string	`json:"first_name,omitempty"`
	LastName		string	`json:"last_name,omitempty"`
	UserImg			string	`json:"image"`
	Status			int		`json:"status"`
}

type FollowerRequest struct{
	TargetId	string	`json:"target_id"`
	Status		int		`json:"status"`
}

func FollowerRequestToEntity(loginUserId string, target FollowerRequest) entity.Follower{
	return entity.Follower{
		SourceId: loginUserId,
		TargetId: target.TargetId,
		Status: target.Status,
	}
}