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

type GroupController struct {
	muxHandler.Handler
	serve  service.GroupServe
	notify service.NotificationServe
}

func GroupHandler(ur service.GroupServe, notify service.NotificationServe) *GroupController {
	gc := &GroupController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		serve:  ur,
		notify: notify,
	}
	gc.InitRoutes()
	return gc
}

func (c *GroupController) InitRoutes() {
	//create new group
	c.Mux.Handle("/new", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var group dto.GroupRequest
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&group); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			groupEntity := dto.GroupRequestToEntity(loggedInUserId, group)
			newGroup, err := c.serve.AddNewGroup(groupEntity)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}

			if err = json.NewEncoder(w).Encode(newGroup); err != nil {
				errHandler.HandleError(w, err)
				return
			}
		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))

	//get all groups
	c.Mux.Handle("/all", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			list, err := c.serve.GetAllGroups()
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

	//get USER CREATED groups (by loggedInUser)
	c.Mux.Handle("/mycreated", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			list, err := c.serve.GetAllMyCreatedGroups(loggedInUserId)
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

	//get USER JOINED groups (by loggedInUser)
	c.Mux.Handle("/joined", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			list, err := c.serve.GetAllMyJoinedGroups(loggedInUserId)
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

	//get 1 POST&Comments or add new POSTS/comments to group
	c.Mux.Handle("/post", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID

		switch r.Method {
		//http://localhost:8080/group/post?groupId=[number]&postId=[number]
		//get all comments of one post
		case "GET":
			keys, ok := r.URL.Query()["groupId"]
			keys2, ok2 := r.URL.Query()["postId"]
		
			if !ok || !ok2 || len(keys[0]) < 1 || len(keys2[0]) < 1 {
				//"id parameter is missing
				logger.WarningLogger.Println("ERROR: no groupId provided")
				e := errHandler.InvalidArgumentError(nil, "ERROR: no groupId provided")
				errHandler.HandleError(w, e)
				return
			} else {
				groupId, err1 := strconv.Atoi(keys[0])
				postId, err2 := strconv.Atoi(keys2[0])
				if err1 != nil {
					logger.WarningLogger.Println("Bad group ID provided:")
					e := errHandler.InvalidArgumentError(err1, "Bad group ID provided: "+err1.Error())
					errHandler.HandleError(w, e)
					return
				}
				if err2 != nil {
					logger.WarningLogger.Println("Bad post ID provided:")
					e := errHandler.InvalidArgumentError(err1, "Bad post ID provided: "+err1.Error())
					errHandler.HandleError(w, e)
					return
				}
				list, err := c.serve.GetOnePostAndComments(loggedInUserId, groupId, postId)
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
			}

		//add new post or comment to particular post in selected group
		case "POST":
			var post dto.GroupPost
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&post); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			postEntity := dto.GroupPostToEntity(loggedInUserId, post)
			err := c.serve.AddNewGroupPost(postEntity)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			fmt.Fprint(w, "NEW POST CREATED IN GROUP SUCCESSFULLY")

		default:
			fmt.Fprint(w, "ONLY GET AND POST METHODS ARE SUPPORTED")
		}
	}))

	//get all posts in one group
	c.Mux.Handle("/post/all", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var groupId int
			var err error
			keys, ok := r.URL.Query()["groupId"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing
				logger.WarningLogger.Println("ERROR: no groupId provided")
				e := errHandler.InvalidArgumentError(nil, "ERROR: no groupId provided")
				errHandler.HandleError(w, e)
				return
			} else {
				groupId, err = strconv.Atoi(keys[0])
				if err != nil {
					logger.WarningLogger.Println("Bad ID provided:", err)
					e := errHandler.InvalidArgumentError(err, "Bad ID provided: "+err.Error())
					errHandler.HandleError(w, e)
					return
				}
			}

			list, err := c.serve.GetOneGroupAllPosts(loggedInUserId, groupId)
			if err != nil {
				if strings.Contains("no group with such ID", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "no group with such ID")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	//GET FRIEND LIST FILTERED “NOT IN PARTICULAR GROUP MEMBERS”
	c.Mux.Handle("/invite/available", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var groupId int
			var err error
			keys, ok := r.URL.Query()["groupId"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing
				logger.WarningLogger.Println("ERROR: no groupId provided")
				e := errHandler.InvalidArgumentError(nil, "ERROR: no groupId provided")
				errHandler.HandleError(w, e)
				return
			} else {
				groupId, err = strconv.Atoi(keys[0])
				if err != nil {
					logger.WarningLogger.Println("Bad ID provided:", err)
					e := errHandler.InvalidArgumentError(err, "Bad ID provided: "+err.Error())
					errHandler.HandleError(w, e)
					return
				}
			}

			result, err := c.serve.GetFriendsNotInGroup(loggedInUserId, groupId)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "you are not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(result); err != nil {
				errHandler.HandleError(w, err)
				return
			}

		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	//invite user to group
	c.Mux.Handle("/invite", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var invitation dto.GroupInvitation
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&invitation); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			result, err := c.serve.InviteUserToGroup(loggedInUserId, invitation)
			if err != nil {
				if strings.Contains("loggedInUser is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "you are not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("user has already invitation", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user has already invitation")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(result); err != nil {
				errHandler.HandleError(w, err)
				return
			}
			// send notification to invited user
			var targetList []string
			targetList = append(targetList, invitation.TargetId)
			err = c.notify.AddNotification(loggedInUserId, targetList, 1, invitation.GroupId)
			if err != nil {
				return
			}

		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))

	//reply to group invitation
	c.Mux.Handle("/invite/reply", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var invitation dto.GroupInvitationReply
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&invitation); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			if invitation.Status > 3 || invitation.Status < 0 {
				logger.WarningLogger.Println("Invalid status received from client:")
				e := errHandler.InvalidArgumentError(errors.New(""), "Invalid status received from client: ")
				errHandler.HandleError(w, e)
				return
			}

			err := c.serve.ReplyToGroupInvitation(loggedInUserId, invitation)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}

			//delete invitation notification from user notification list
			err3 := c.notify.DeleteNotification(invitation.ActorId, loggedInUserId, 1, invitation.GroupId)
			if err3 != nil {
				logger.ErrorLogger.Println(err)
			}

			fmt.Fprint(w, "Reply was added to DB")
		default:
			fmt.Fprint(w, "ONLY PUT METHOD IS SUPPORTED")
		}
	}))

	//make join-group request
	c.Mux.Handle("/join", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID

		switch r.Method {
		case "POST":
			var request dto.GroupId
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&request); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			newRequest, result, err := c.serve.RequestGroupAccess(loggedInUserId, request.GroupId)
			if err != nil {
				if strings.Contains("user already have invitation", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user already have invitation")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if err = json.NewEncoder(w).Encode(result); err != nil {
				errHandler.HandleError(w, err)
				return
			}
			// send notification to group admin
			if newRequest{
				admin, err := c.serve.GetGroupAdmin(request.GroupId)
				if err != nil {
					logger.WarningLogger.Println("ERROR:", err)
				}
				var targetList []string
				targetList = append(targetList, admin)
				err = c.notify.AddNotification(loggedInUserId, targetList, 2, request.GroupId)
				if err != nil {
					return
				}
			}
			
		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))

	// pending join-group request (to admin only)
	c.Mux.Handle("/join/reply", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)
		loggedInUserId := val.UserID
		switch r.Method {
		case "GET":
			// get pending join-group request (to admin only)
			keys, ok := r.URL.Query()["groupId"]
			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing
				logger.WarningLogger.Println("ERROR: no groupId provided")
				e := errHandler.InvalidArgumentError(nil, "ERROR: no groupId provided")
				errHandler.HandleError(w, e)
				return
			} else {
				groupId, err := strconv.Atoi(keys[0])
				if err != nil {
					logger.WarningLogger.Println("ERROR: BAD groupId provided")
					e := errHandler.InvalidArgumentError(nil, "ERROR: BAD groupId provided")
					errHandler.HandleError(w, e)
					return
				}
				list, err := c.serve.GetPendingJoinRequests(loggedInUserId, groupId)
				if err != nil {
					if strings.Contains("loggedInUser is not an admin in group", err.Error()) {
						foundError := errHandler.InvalidArgumentError(nil, "not enough rights")
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
			// reply to join-group request (to admin only)
			var reply dto.GroupAccessRequestReply
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&reply); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			if reply.Status < 0 || reply.Status > 3 {
				logger.WarningLogger.Println("Invalid STATUS received from client")
				e := errHandler.InvalidArgumentError(errors.New(""), "Invalid STATUS received from client ")
				errHandler.HandleError(w, e)
				return
			}
			err := c.serve.ReplyToGroupAccessRequest(loggedInUserId, reply)
			if err != nil {
				if strings.Contains("user already have invitation", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user already have invitation")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			if reply.Status == 1 {
				// send notification to access_requested_user
				var targetList []string
				targetList = append(targetList, reply.TargetId)
				err = c.notify.AddNotification(loggedInUserId, targetList, 7, reply.GroupId)
				if err != nil {
					return
				}
			}

			//delete the notification from admin notification list
			err3 := c.notify.DeleteNotification(reply.TargetId, loggedInUserId, 2, reply.GroupId)
			if err3 != nil {
				logger.ErrorLogger.Println(err)
			}

			fmt.Fprint(w, "REQUEST WAS UPDATED")
		default:
			fmt.Fprint(w, "ONLY GET AND PUT METHODS ARE SUPPORTED")
		}
	}))

	//create new group event (any group member)
	c.Mux.Handle("/event/new", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var event dto.GroupEvent
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&event); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			entity := dto.GroupEventToEntity(loggedInUserId, event)

			eventId, err := c.serve.CreateEvent(entity)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			//send notification to all group members about new event
			groupList, err := c.serve.GetAllGroupMembersExceptMe(loggedInUserId, eventId)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
			}
			err = c.notify.AddNotification(loggedInUserId, groupList, 3, eventId)
			if err != nil {
				logger.WarningLogger.Println("ERROR:", err)
			}

			fmt.Fprint(w, "NEW EVENT CREATED")
		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))

	//group member can reply to event with going/not going status
	c.Mux.Handle("/event/reply", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			var event dto.EventParticipant
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields()
			if err := d.Decode(&event); err != nil {
				logger.WarningLogger.Println("Invalid json received from client:", err)
				e := errHandler.InvalidArgumentError(err, "Invalid json received from client: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			err := c.serve.AddEventParticipant(loggedInUserId, event)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("no such event or group presented", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "no such event or group presented")
					errHandler.HandleError(w, foundError)
					return
				}
				logger.WarningLogger.Println("ERROR:", err)
				e := errHandler.DataBaseError(err)
				errHandler.HandleError(w, e)
				return
			}
			fmt.Fprint(w, "USER REPLY ADDED")
		default:
			fmt.Fprint(w, "ONLY POST METHOD IS SUPPORTED")
		}
	}))

	//get all group events
	c.Mux.Handle("/event/all", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			keys, ok := r.URL.Query()["groupId"]

			if !ok || len(keys[0]) < 1 {
				//"id parameter is missing
				logger.WarningLogger.Println("ERROR: no groupId provided")
				e := errHandler.InvalidArgumentError(nil, "ERROR: no groupId provided")
				errHandler.HandleError(w, e)
				return
			} else {
				groupId, err := strconv.Atoi(keys[0])
				if err != nil {
					logger.WarningLogger.Println("ERROR: BAD groupId provided")
					e := errHandler.InvalidArgumentError(nil, "ERROR: BAD groupId provided")
					errHandler.HandleError(w, e)
					return
				}

				list, err := c.serve.GetAllGroupEvents(loggedInUserId, groupId)
				if err != nil {
					if strings.Contains("user is not a group member", err.Error()) {
						foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
						errHandler.HandleError(w, foundError)
						return
					}
					if strings.Contains("no such group presented", err.Error()) {
						foundError := errHandler.InvalidArgumentError(nil, "no such group presented")
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	// get all (FUTURE) events CREATED BY ME
	c.Mux.Handle("/event/mycreated", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			list, err := c.serve.GetMyCreatedEvents(loggedInUserId)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("no such group presented", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "no such group presented")
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	// get all (FUTURE) events GOING BY ME / I participate
	c.Mux.Handle("/event/joined", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			list, err := c.serve.GetMyJoinedEvents(loggedInUserId)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("no such group presented", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "no such group presented")
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	//GET ALL MESSAGES IN ONE GROUP CHAT
	c.Mux.Handle("/chat", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, _ := r.Context().Value("values").(middleware.UserContext)
			sender := val.UserID
			receiver := r.URL.Query().Get("groupId")
			skip := r.URL.Query().Get("skip")
			limit := r.URL.Query().Get("limit")
			intSkip, _ := strconv.Atoi(skip)
			intLimit, _ := strconv.Atoi(limit)

			groupId, err := strconv.Atoi(receiver)
			if err != nil {
				logger.WarningLogger.Println("Bad Group ID provided:", err)
				e := errHandler.InvalidArgumentError(err, "Bad group ID provided: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}
			

			list, err := c.serve.GetOneGroupAllMessages(sender, groupId, intSkip, intLimit)
			if err != nil {
				if strings.Contains("user is not a group member", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "user is not a group member")
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}
	}))

	//get one group information (by groupId in path parameter)
	c.Mux.Handle("/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			//get one group info (description)
			val, _ := r.Context().Value("values").(middleware.UserContext)
			loggedInUserId := val.UserID

			stringID := strings.TrimPrefix(r.URL.Path, "/")
			groupId, err := strconv.Atoi(stringID)
			if err != nil {
				logger.WarningLogger.Println("Bad ID provided:", err)
				e := errHandler.InvalidArgumentError(err, "Bad ID provided: "+err.Error())
				errHandler.HandleError(w, e)
				return
			}

			list, err := c.serve.GetOneGroupInfo(loggedInUserId, groupId)
			if err != nil {
				if strings.Contains("no group with such ID", err.Error()) {
					foundError := errHandler.InvalidArgumentError(nil, "no group with such ID")
					errHandler.HandleError(w, foundError)
					return
				}
				if strings.Contains("user is not a group member", err.Error()) {
					shortList := dto.GroupShortInfo{
						Id:          list.Id,
						Title:       list.Title,
						Description: list.Description,
						Members:     len(list.Members),
					}

					if err = json.NewEncoder(w).Encode(shortList); err != nil {
						errHandler.HandleError(w, err)
						return
					}
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
		default:
			fmt.Fprint(w, "ONLY GET METHOD IS SUPPORTED")
		}

	}))
}
