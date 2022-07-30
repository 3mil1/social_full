package entity

type User struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	PwHash      string `json:"pw_hash,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	AboutMe     string `json:"about_me,omitempty"`
	UserImg     string `json:"user_image"`
	IsPrivate   bool   `json:"is_private,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

//type Claims struct {
//	ID        string `json:"ID"`
//	FirstName string `json:"first_name,omitempty"`
//	LastName  string `json:"last_name,omitempty"`
//	jwt.MapClaims
//}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type RefreshTokenDB struct {
	RefreshToken string
	UserID       string
	Device       string
	IP           string
}
