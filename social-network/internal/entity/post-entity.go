package entity

type UserPost struct{
	Id			int		 	`json:"id"`
	UserId		string	 	`json:"user_id"`
	Subject		string		`json:"subject"`
	Content		string		`json:"content"`
	Image		string		`json:"image"`
	ParentId	int			`json:"parent_id"`
	CreatedAt	string		`json:"created_at"`
	UpdatedAt	string		`json:"updated_at"`
	Privacy		int			`json:"privacy"`
	Access 		[]string	`json:"access"`
}

type PostAccess struct{
	UserId		string
	FirstName	string
	LastName	string
}

