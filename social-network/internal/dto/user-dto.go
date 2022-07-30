package dto

import "social-network/internal/entity"

type UserRequestBody struct {
	Email     string `json:"email" validate:"required,max=30,min=2"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,max=30,min=2"`
	LastName  string `json:"last_name" validate:"required,max=30,min=2"`
	BirthDay  string `json:"birth_day" validate:"required"`
	Nickname  string `json:"nickname" validate:"max=10"`
	AboutMe   string `json:"about_me" validate:"max=30"`
	UserImg   string `json:"user_img"`
}

func (user UserRequestBody) ToUserEntity() *entity.User {
	return &entity.User{
		Email:       user.Email,
		PwHash:      user.Password,
		LastName:    user.LastName,
		FirstName:   user.FirstName,
		DateOfBirth: user.BirthDay,
		Nickname:    user.Nickname,
		AboutMe:     user.AboutMe,
		UserImg:     user.UserImg,
	}
}

type UserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDay  string `json:"birth_day"`
	Nickname  string `json:"nickname"`
	AboutMe   string `json:"about_me"`
	UserImg   string `json:"user_img"`
	IsPrivate bool   `json:"is_private"`
}

func CreateUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDay:  user.DateOfBirth,
		Nickname:  user.Nickname,
		AboutMe:   user.AboutMe,
		UserImg:   user.UserImg,
		IsPrivate: user.IsPrivate,
	}
}

type SignInRequestBody struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	ID        string `json:"ID,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func ToRefreshTokenDB(tr RefreshTokenRequestBody) entity.RefreshTokenDB {
	return entity.RefreshTokenDB{
		RefreshToken: tr.RefreshToken,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserProfileRequest struct {
	ID string `json:"ID,omitempty"`
}

type PrivateProfileResponse struct {
	ID        string `json:"ID,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserImg   string `json:"user_img"`
}

type UserUpdate struct {
	Nickname  string `json:"nickname" validate:"max=10"`
	AboutMe   string `json:"about_me" validate:"max=30"`
	UserImg   string `json:"user_img"`
	IsPrivate bool   `json:"is_private"`
}
