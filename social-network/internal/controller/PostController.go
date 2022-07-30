package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social-network/internal/app/middleware"
	"social-network/internal/dto"
	"social-network/internal/service"
	"social-network/pkg/errHandler"
	"social-network/pkg/logger"
	"social-network/pkg/muxHandler"
	"strconv"
	"strings"
)

type PostController struct {
	muxHandler.Handler
	serve service.PostServe
	notify service.NotificationServe
}

func PostHandler(fr service.PostServe, notify service.NotificationServe) *PostController {
	pc := &PostController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		serve: fr,
		notify: notify,
	}
	pc.InitPostRoutes()
	return pc
}

func (c *PostController) InitPostRoutes() {
	c.Mux.Handle("/new", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var post dto.PostReceive
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&post); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			if post.Privacy > 3 || post.Privacy < 1{
				logger.WarningLogger.Println("Invalid post privacy received from client")
				e := errHandler.InvalidArgumentError(errors.New(""), "Invalid post privacy received from client")
				errHandler.HandleError(w, e)
				return
			}
			post.UserId = loggedInUserId
			result, err := c.serve.AddNewPost(post)
			if err != nil {
				logger.WarningLogger.Println("ERROR: AddNewPost:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(result); err != nil {
				errHandler.HandleError(w, err)
				return
			}
			//send notification to post_owner, if comment is added
			if post.ParentId > 0 {
				postOwner, err3 := c.serve.GetPostOwner(post.ParentId)
				if err3 != nil{
					return
				}
				if postOwner != loggedInUserId{
					var targetList []string
					targetList = append(targetList, postOwner)
					err_new := c.notify.AddNotification(loggedInUserId, targetList, 6, post.ParentId)
					if err_new != nil {
						logger.ErrorLogger.Println(err)
					}
				}
			}

		default:
			fmt.Fprintf(w, "Sorry, only POST method is supported.")
		}
	}))
	c.Mux.Handle("/oneuser", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID
			var userId string
			keys, ok := r.URL.Query()["id"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing. User want`s to get his own posts
				userId = loggedInUserId
			} else {
				userId = keys[0]
			}
			list, err := c.serve.GetAllUserPosts(loggedInUserId, userId)
			if err != nil {
				logger.WarningLogger.Println("ERROR: GetAllUserPosts:", err)
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

	c.Mux.Handle("/all", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			list, err := c.serve.GetAllPosts(loggedInUserId)
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
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	c.Mux.Handle("/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		stringID := strings.TrimPrefix(r.URL.Path, "/")
		postId, err := strconv.Atoi(stringID)
		if err != nil {
			logger.WarningLogger.Println("Bad ID provided:", err)
			e := errHandler.InvalidArgumentError(err, "Bad ID provided: "+err.Error())
			errHandler.HandleError(w, e)
			return
		}
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID

		list, err := c.serve.GetOnePostWithComments(loggedInUserId, postId)
		if err != nil {
			if strings.Contains("user has no access to this post", err.Error()) {
				foundError := errHandler.InvalidArgumentError(nil, "user has no access to this post")
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
	}))
}
