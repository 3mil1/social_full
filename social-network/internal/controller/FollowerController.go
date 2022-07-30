package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social-network/internal/app/middleware"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/internal/service"
	"social-network/pkg/errHandler"
	"social-network/pkg/logger"
	"social-network/pkg/muxHandler"
	"strings"
)

type FollowerController struct {
	muxHandler.Handler
	serve  service.FollowerServe
	notify service.NotificationServe
}

func FollowerHandler(fr service.FollowerServe, notify service.NotificationServe) *FollowerController {
	fc := &FollowerController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		serve:  fr,
		notify: notify,
	}
	fc.InitFollowerRoutes()
	return fc
}

func (contr *FollowerController) InitFollowerRoutes() {
	contr.Mux.Handle("/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID

		switch r.Method {
			// GET ALL users  loggedInUser is following
		case "GET":
			var userId string
			keys, ok := r.URL.Query()["id"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing. User want`s to get his own followers
				userId = loggedInUserId
			} else {
				userId = keys[0]
			}
			list, err := contr.serve.GetAllUsersIFollow(loggedInUserId, userId)
			if err != nil {
				logger.WarningLogger.Println("ERROR: GetAllUsersIFollow:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(list); err != nil {
				errHandler.HandleError(w, err)
				return
			}
		case "POST":
			//ADD NEW FOLLOWER
			var followerRequest dto.FollowerRequest
			err := json.NewDecoder(r.Body).Decode(&followerRequest)
			if err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			follower := dto.FollowerRequestToEntity(loggedInUserId, followerRequest)
			connection, err := contr.serve.AddNewFollower(follower)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			} else if connection == 1 { //only if new connection added -> send notification
				var targetList []string
				targetList = append(targetList, follower.TargetId)
				err = contr.notify.AddNotification(follower.SourceId, targetList, 4, 0)
				if err != nil {
					return
				}
				fmt.Fprintf(w, "NEW CONNECTION ADDED TO DB")
			}
		case "PUT":
			//reply to follower request
			var followerRequest dto.FollowerRequest
			err := json.NewDecoder(r.Body).Decode(&followerRequest)
			if err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			if followerRequest.Status > 2 {
				logger.WarningLogger.Println("Invalid status received from client")
				e := errHandler.InvalidArgumentError(err, "Invalid status received from client")
				errHandler.HandleError(w, e)
				return
			}
			follower := entity.Follower{SourceId: followerRequest.TargetId, TargetId: loggedInUserId, Status: followerRequest.Status}
			err = contr.serve.UpdateFollower(follower)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			} else {
				//delete the notification for target_user
				err3 := contr.notify.DeleteNotification(followerRequest.TargetId, loggedInUserId, 4, 0)
				if err3 != nil {
					logger.ErrorLogger.Println(err)
				}
				fmt.Fprintf(w, "CONNECTION UPDATED")
			}
		case "DELETE":
			//user wants to unfollow
			var userId string
			keys, ok := r.URL.Query()["id"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing. User want`s to get his own followers
				logger.WarningLogger.Println("ERROR: ID PARAMETER MISSING")
				e := errHandler.InvalidArgumentError(errors.New("NO ID"), "ID PARAMETER IS MISSING")
				errHandler.HandleError(w, e)
				return
			} else {
				userId = keys[0]
			}
			err := contr.serve.DeleteFollower(loggedInUserId, userId)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			//check notofication table and if notification about friend request exist delete it
			err_new := contr.notify.DeleteNotification(loggedInUserId, userId, 4, 0)
			if err_new != nil{
				logger.ErrorLogger.Println(err)
			}

			fmt.Fprintf(w, "User follower request/connection deleted")
		default:
			fmt.Fprintf(w, "Sorry, only GET, POST, PUT  methods are supported.")
		}
	}))

	contr.Mux.Handle("/back", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID
		switch r.Method {
		case "GET":
			var userId string
			keys, ok := r.URL.Query()["id"]
			if !ok || len(keys[0]) < 1 {
				/*if loggedInUser request his own information					*/
				list, err := contr.serve.GetAllUsersFollowsMe(loggedInUserId)
				if err != nil {
					logger.WarningLogger.Println("ERROR: GetAllUsersFollowsMe:", err)
					e := errHandler.DataBaseError(err)
					errHandler.HandleError(w, e)
					return
				}
				if err = json.NewEncoder(w).Encode(list); err != nil {
					errHandler.HandleError(w, err)
					return
				}
			} else {
				userId = keys[0]
				/*if loggedInUser request some user information					*/

				list, err := contr.serve.GetUsersFollowsMeAcceptedStatusOnly(loggedInUserId, userId)
				if err != nil {
					if strings.Contains("you are not connected", err.Error()) {
						foundError := errHandler.InvalidArgumentError(nil, "user has no access to this user")
						errHandler.HandleError(w, foundError)
						return
					}
					logger.WarningLogger.Println("ERROR:", err)
					e := errHandler.DataBaseError(err)
					errHandler.HandleError(w, e)
					return
				}
				if err = json.NewEncoder(w).Encode(list); err != nil {
					errHandler.HandleError(w, err)
					return
				}
			}
		case "PUT":
			var followerRequest dto.FollowerRequest
			err := json.NewDecoder(r.Body).Decode(&followerRequest)
			if err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			if followerRequest.Status > 2 {
				logger.WarningLogger.Println("Invalid status received from client")
				e := errHandler.InvalidArgumentError(err, "Invalid status received from client")
				errHandler.HandleError(w, e)
				return
			}
			follower := dto.FollowerRequestToEntity(loggedInUserId, followerRequest)
			err = contr.serve.UpdateFollowRequest(follower)
			if err != nil {
				logger.WarningLogger.Println("ERROR: UpdateFollowRequest:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			} else {
				fmt.Fprintf(w, "CONNECTION UPDATED")
			}
		default:
			fmt.Fprintf(w, "Sorry, only GET and PUT methods are supported.")
		}
	}))
	contr.Mux.Handle("/chat", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID
		switch r.Method {
		case "GET":
			//get a list of users with whom I can chat
			list, err := contr.serve.GetChatList(loggedInUserId)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
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
}
