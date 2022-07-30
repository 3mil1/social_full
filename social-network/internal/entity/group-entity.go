package entity

type Group struct {
	Id          int
	Creator     string
	Title       string
	Description string
}

type GroupPost struct{
	Id 			int
	GroupId		int
	CreatorId		string
	Title		string
	Content		string
	Image		string
	ParentId	int
	CreatedAt	int
	UpdatedAt	int
}

type GroupEventEntity struct{
	Id 			int
	GroupId 	int
	UserId 		string
	Title 		string
	Description string
	EventDate 	string
	CreatedAt 	string
	GoingStatus int
}

type GroupMessage struct{
	Id 				int
	UserId 		string
	GroupId 		int		
	Content			string
	CreatedAt 		string
}