package service

import (
	"errors"
	"fmt"
	"os"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/env"
	"social-network/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{
		r,
	}
}

func (service *UserService) AddUser(data dto.UserRequestBody) (*dto.UserResponse, error) {
	user := data.ToUserEntity()
	user.ID = uuid.New().String()
	hash, _ := HashPassword(user.PwHash)
	user.PwHash = hash
	s := strings.Split(user.DateOfBirth, "-")
	day, _ := strconv.Atoi(s[0])
	month, _ := strconv.Atoi(s[1])
	year, _ := strconv.Atoi(s[2])

	// convert from 01-02-2006 time to unix
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	user.DateOfBirth = strconv.FormatInt(t.Unix(), 10)

	userResp, err := service.repo.AddUser(user)
	if err != nil {
		return nil, err
	}
	// convert from unix time to 01-02-2006
	i, err := strconv.ParseInt(userResp.DateOfBirth, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	userResp.DateOfBirth = tm.Format("2006-01-02") //	yyyy-MM-dd
	newUser := dto.CreateUserResponse(userResp)

	return &newUser, nil
}


func (service *UserService) SignIn(data dto.SignInRequestBody, ip, userAgent string) (*dto.TokenResponse, error) {
	user, err := service.repo.GetUserByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(user.PwHash, data.Password) {
		logger.ErrorLogger.Println(err)
		return nil, errors.New("wrong pw")
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := &entity.RefreshTokenDB{
		UserID:       user.ID,
		RefreshToken: token.RefreshToken,
		Device:       userAgent,
		IP:           ip,
	}

	err = service.repo.AddRefreshToken(*refreshToken)
	if err != nil {
		return nil, err
	}

	response := &dto.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return response, nil
}

func CreateToken(userId string) (*entity.TokenDetails, error) {
	td := &entity.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 500).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()

	var err error
	//Creating Access Token
	accessSecret := env.GoDotEnvVariable("ACCESS_SECRET")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv(accessSecret)))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	refreshSecret := env.GoDotEnvVariable("REFRESH_SECRET")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv(refreshSecret)))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (service *UserService) RefreshToken(refreshTokenFromRequest dto.RefreshTokenRequestBody, ip, userAgent string) (*dto.TokenResponse, error) {
	//verify the refresh_token
	refreshSecret := env.GoDotEnvVariable("REFRESH_SECRET")
	token, err := jwt.Parse(refreshTokenFromRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the refresh_token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(refreshSecret)), nil
	})
	//if there is an error, the refresh_token must have expired
	if err != nil {
		logger.WarningLogger.Println("Refresh refresh_token expired")
		return nil, err
	}
	//is refresh_token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		logger.WarningLogger.Println(err)
		return nil, err
	}
	//Since refresh_token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the refresh_token claims should conform to MapClaims
	if ok && token.Valid {
		_, ok = claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			logger.WarningLogger.Println(err)
			return nil, err
		}

		//get uID from refresh token
		userId := fmt.Sprintf("%s", claims["user_id"])
		rt := dto.ToRefreshTokenDB(refreshTokenFromRequest)
		rt.UserID = userId

		rtFromDB, err := service.repo.DeleteRefreshToken(rt)
		if err != nil {
			return nil, err
		}

		if !(rtFromDB.IP == ip && rtFromDB.Device == userAgent) {
			logger.WarningLogger.Println("Device or Ip is different")
			return nil, errors.New("please login again")
		}

		//Create new pairs of refresh and access tokens
		ts, err := CreateToken(userId)
		if err != nil {
			logger.WarningLogger.Println(err)
			return nil, err
		}

		refreshToken := &entity.RefreshTokenDB{
			UserID:       userId,
			RefreshToken: ts.RefreshToken,
			Device:       userAgent,
			IP:           ip,
		}

		err = service.repo.AddRefreshToken(*refreshToken)
		if err != nil {
			return nil, err
		}

		var tokens = &dto.TokenResponse{}
		tokens.AccessToken = ts.AccessToken
		tokens.RefreshToken = ts.RefreshToken

		return tokens, nil
	}
	return nil, err
}

func (service *UserService) GetUserByID(userID string) (*dto.UserResponse, error) {
	user, err := service.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	newUser := dto.CreateUserResponse(user)
	return &newUser, nil
}

func (service *UserService) SignOut(userID, ip, userAgent string) error {
	err := service.repo.DeleteSession(userID, ip, userAgent)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (service *UserService) UpdateUser(user dto.UserUpdate, loggedInUser string) error {
	return service.repo.UpdateUser(user, loggedInUser)

}

func (service *UserService) GetMyFollowerProfile(loggedInUserId, requestedUserId string) (*dto.UserResponse, error) {
	/*check that requested user is PUBLIC profile
	TRUE - PRIVATE
	FALSE - PUBLIC
	*/
	var user dto.UserResponse
	private, err := service.repo.GetUserStatusByID(requestedUserId)
	if err != nil {
		return &user, err
	}
	if !private {
		//status is public-> send FULL profile
		return service.GetUserByID(requestedUserId)
	} else {
		//IF requested user is PRIVATE profile -> check followers connection
		connection, err := service.repo.Get2UsersConnectionStatus(loggedInUserId, requestedUserId)
		if err != nil {
			return &user, err
		}
		if connection == 1 {
			return service.GetUserByID(requestedUserId)
		} else if connection != 1 {
			user, err := service.GetUserByID(requestedUserId)
			if err != nil {
				return user, nil
			}
			user.Email = ""
			user.BirthDay = ""
			user.AboutMe = ""
			return user, nil
		}
	}

	return &user, nil
}

func (service *UserService) GetAllUsers() ([]dto.PrivateProfileResponse, error) {
	return service.repo.GetAllUsers()
}
