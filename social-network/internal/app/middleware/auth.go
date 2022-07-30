package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"social-network/pkg/env"
	"social-network/pkg/errHandler"
	"strings"
)

type UserContext struct {
	UserID string
}

func Auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := TokenValid(r)
		if err != nil {
			er := errHandler.Unauthorized(err)
			errHandler.HandleError(w, er)
			return
		}

		ctx := context.WithValue(r.Context(), "values", UserContext{UserID: fmt.Sprintf("%s", token.Claims.(jwt.MapClaims)["user_id"])})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenValid(r *http.Request) (*jwt.Token, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}
	return token, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	accessSecret := env.GoDotEnvVariable("ACCESS_SECRET")
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the refresh_token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(accessSecret)), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func GetIP(r *http.Request) []string {
	ip := make([]string, 0)
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		ip = append(ip, forwarded)
		return ip
	}
	ip = append(ip, r.RemoteAddr)
	return ip
}
func UserAgent(r *http.Request) []string {
	userAgent := r.Header["User-Agent"]
	return userAgent
}
