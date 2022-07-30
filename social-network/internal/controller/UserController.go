package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"social-network/internal/app/middleware"
	"social-network/internal/dto"
	"social-network/internal/service"
	"social-network/pkg/errHandler"
	"social-network/pkg/logger"
	"social-network/pkg/muxHandler"
	"social-network/pkg/validate"
	"strings"
)

type UserController struct {
	muxHandler.Handler
	u service.User
	notify service.NotificationServe
	chat service.ChatServe
}

func UsersHandler(ur service.User, notify service.NotificationServe, chat service.ChatServe) *UserController {
	uh := &UserController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		u: ur,
		notify: notify,
		chat: chat,
	}
	uh.InitRoutes()
	return uh
}

func (c *UserController) InitRoutes() {
	c.Mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var user dto.UserRequestBody

		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&user); err != nil {
			logger.WarningLogger.Println("Invalid json received from client:", err)
			e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
			errHandler.HandleError(w, e)
			return
		}

		//validate
		e := validate.Struct(&dto.UserRequestBody{}, &user)
		if reflect.ValueOf(e).Len() > 0 {
			if err := json.NewEncoder(w).Encode(e); err != nil {
				logger.ErrorLogger.Println(err)
				return
			}
			return
		}

		createUser, err := c.u.AddUser(user)
		if err != nil {
			if strings.Contains(fmt.Sprintf("%s", err.Error()), "UNIQUE") {
				foundError := errHandler.InvalidArgumentError(nil, "user with this email already exists")
				errHandler.HandleError(w, foundError)
				return
			}
			errHandler.HandleError(w, err)
			return
		}

		logger.InfoLogger.Println("NewUserService user was added to DB")

		if err = json.NewEncoder(w).Encode(createUser); err != nil {
			errHandler.HandleError(w, err)
			return
		}

		return
	})

	c.Mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		var ip = middleware.GetIP(r)
		var userAgent = middleware.UserAgent(r)
		var loginReq = &dto.SignInRequestBody{}

		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&loginReq); err != nil {
			logger.WarningLogger.Println("Invalid json received from client:", err)
			e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
			errHandler.HandleError(w, e)
			return
		}

		logger.InfoLogger.Println(ip, userAgent)
		//from 172.19.0.1:61434 to 172.19.0.1
		ipSplit := strings.Split(ip[0], ":")
		jwt, err := c.u.SignIn(*loginReq, ipSplit[0], userAgent[0])
		if err != nil {
			if strings.Contains(fmt.Sprintf("%s", err.Error()), "no rows in result set") || strings.Contains(fmt.Sprintf("%s", err.Error()), "wrong pw") {
				foundError := errHandler.InvalidArgumentError(nil, "password or email is incorrect")
				errHandler.HandleError(w, foundError)
				return
			}
			if strings.Contains(fmt.Sprintf("%s", err.Error()), "no rows in result set") {
				foundError := errHandler.InvalidArgumentError(nil, "password or email is incorrect")
				errHandler.HandleError(w, foundError)
				return
			}
			errHandler.HandleError(w, err)
			return
		}

		if err = json.NewEncoder(w).Encode(jwt); err != nil {
			errHandler.HandleError(w, err)
			return
		}
	})

	c.Mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		var ip = middleware.GetIP(r)
		var userAgent = middleware.UserAgent(r)

		var refreshToken dto.RefreshTokenRequestBody

		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&refreshToken); err != nil {
			logger.WarningLogger.Println("Invalid json received from client:", err)
			e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
			errHandler.HandleError(w, e)
			return
		}

		logger.InfoLogger.Println(ip, userAgent)
		//from 172.19.0.1:61434 to 172.19.0.1
		ipSplit := strings.Split(ip[0], ":")
		newTokenPair, err := c.u.RefreshToken(refreshToken, ipSplit[0], userAgent[0])
		if err != nil {
			e := errHandler.LoginAgain(err)
			errHandler.HandleError(w, e)
			return
		}

		if err = json.NewEncoder(w).Encode(newTokenPair); err != nil {
			errHandler.HandleError(w, err)
			return
		}
	})

	c.Mux.Handle("/me", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			user, err := c.u.GetUserByID(val.UserID)
			if err != nil {
				logger.WarningLogger.Println("ERROR: GetUserByID:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}

			if err = json.NewEncoder(w).Encode(user); err != nil {
				errHandler.HandleError(w, err)
				return
			}
		case "PUT":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUser := val.UserID

			var user dto.UserUpdate

			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&user); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			//validate
			e := validate.Struct(&dto.UserUpdate{}, &user)
			if reflect.ValueOf(e).Len() > 0 {
				if err := json.NewEncoder(w).Encode(e); err != nil {
					logger.ErrorLogger.Println(err)
					return
				}
				return
			}

			err := c.u.UpdateUser(user, loggedInUser)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			fmt.Fprintf(w, "User is updated")
		default:
			fmt.Fprintf(w, "Sorry, only GET and PUT methods are supported.")
		}
	}))

	c.Mux.Handle("/signout", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		var ip = middleware.GetIP(r)
		var userAgent = middleware.UserAgent(r)

		val, _ := r.Context().Value("values").(middleware.UserContext)
		//from 172.19.0.1:61434 to 172.19.0.1
		ipSplit := strings.Split(ip[0], ":")
		err := c.u.SignOut(val.UserID, ipSplit[0], userAgent[0])
		if err != nil {
			return
		}
	}))

	c.Mux.Handle("/oneuser", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			keys, ok := r.URL.Query()["id"]

			if !ok || len(keys[0]) < 1 {
				fmt.Fprint(w, "id parameter is missing.")
				return
			}
			// Query()["id"] will return an array of items,
			// we only want the single item.
			userId := keys[0]
			user, err := c.u.GetMyFollowerProfile(loggedInUserId, userId)
			if err != nil {
				logger.WarningLogger.Println("ERROR: GetMyFollowerProfile:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(user); err != nil {
				errHandler.HandleError(w, err)
				return
			}
		default:
			fmt.Fprintf(w, "Sorry, only GET method is supported.")
		}
	}))

	c.Mux.Handle("/all", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			list, err := c.u.GetAllUsers()
			if err != nil {
				logger.WarningLogger.Println("ERROR: GetAllUsers:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(list); err != nil {
				errHandler.HandleError(w, err)
				return
			}
		default:
			fmt.Fprintf(w, "Sorry, only GET method is supported.")
		}
	}))
	c.Mux.Handle("/notification/reply", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID

		switch r.Method {
		case "POST":
			keys, ok := r.URL.Query()["id"]
			keys2 := r.URL.Query()["status"]

			if !ok || len(keys[0]) < 1 {
				fmt.Fprint(w, "id parameter is missing.")
				return
			}
			notify_id := keys[0]
			status := keys2[0]
			if status != "1" && status != "2"{
				fmt.Fprint(w, "wrong status")
				return
			}
	
			err := c.notify.UpdateNotification(loggedInUserId, notify_id, keys2[0])
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			fmt.Fprint(w, "notification seen-status changed")

		//delete notifications for personal and group chat
		case "DELETE":
			keys, ok := r.URL.Query()["id"]
			if !ok || len(keys[0]) < 1 {
				fmt.Fprint(w, "id parameter is missing.")
				return
			}
			target_id := keys[0]
			//mark all messages as seen to loggedInUser
			c.chat.MarkMessageAsSeen(loggedInUserId, target_id)

			//delete notifications from list for LoggedInUser
			if len(target_id) > 20{
				//id is user_id
				err := c.notify.DeleteNotification(target_id, loggedInUserId, 5, 0) 
				if err != nil {
					logger.WarningLogger.Println("ERROR:", err)
					e := errHandler.DataBaseError(err)
					errHandler.HandleError(w, e)
					return
				}
			} else {
				//id is group_id
				err := c.notify.DeleteGroupChatNotification(loggedInUserId, target_id)
				if err != nil {
					logger.WarningLogger.Println("ERROR:", err)
					e := errHandler.DataBaseError(err)
					errHandler.HandleError(w, e)
					return
				}
			}
			fmt.Fprint(w, "notification seen-status changed")
		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))
}
