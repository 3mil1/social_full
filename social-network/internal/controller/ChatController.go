package controller

import (
	"encoding/json"
	"net/http"
	"social-network/internal/app/middleware"
	"social-network/internal/service"
	"social-network/pkg/muxHandler"
	"strconv"
)

type ChatController struct {
	muxHandler.Handler
	s service.ChatServe
}

func ChatHandler(cs service.ChatServe) *ChatController {
	ch := &ChatController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		s: cs,
	}
	ch.InitRoutes()
	return ch
}

func (c *ChatController) InitRoutes() {
	//personal chat
	c.Mux.Handle("/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		val, _ := r.Context().Value("values").(middleware.UserContext)

		sender := val.UserID
		receiver := r.URL.Query().Get("with")
		skip := r.URL.Query().Get("skip")
		limit := r.URL.Query().Get("limit")
		intSkip, _ := strconv.Atoi(skip)
		intLimit, _ := strconv.Atoi(limit)

		messages, err := c.s.GetMessages(sender, receiver, intSkip, intLimit)
		if err != nil {
			//handleError(w, err)
			return
		}

		if err = json.NewEncoder(w).Encode(messages); err != nil {
			//handleError(w, err)
			return
		}

	}))

	//group chat is in group controller
}
