package controller

import (
	"net/http"
	"social-network/internal/dto"
	"social-network/internal/service"
	"social-network/pkg/logger"
	"social-network/pkg/muxHandler"
	"sort"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsController struct {
	muxHandler.Handler
	notify service.NotificationServe
	chat   service.ChatServe
}

func WsHandler(notify service.NotificationServe, chat service.ChatServe) *WsController {
	n := &WsController{
		Handler: muxHandler.Handler{
			Mux: http.NewServeMux(),
		},
		notify: notify,
		chat:   chat,
	}
	n.InitRoutes()
	return n
}

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConn]string)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConn struct {
	*websocket.Conn
}

type WsPayload struct {
	Action         string        `json:"action,omitempty"`
	MessageTo      string        `json:"message_to,omitempty"`	//userID / group_number
	MessageContent string        `json:"message_content,omitempty"`
	UserID         string        `json:"user,omitempty"`
	Conn           WebSocketConn `json:"-"`
}

func (c *WsController) reader(conn *WebSocketConn) {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorLogger.Println(r)
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func (c *WsController) ListenToWsChannel() {
	//var response dto.WsResponse

	for {
		e := <-wsChan
		switch e.Action {
		case "connect":
			clients[e.Conn] = e.UserID
			notifications, err := c.notify.GetAllNotifications(e.UserID)
			if err != nil {
				
				return
			}
			c.SendOne(notifications, e.UserID)
		case "message":
			messages, err := c.chat.AddMessage(e.UserID, e.MessageTo, e.MessageContent)
			if err != nil {
				return
			}
			c.SendOne(messages, e.UserID)
			c.SendOne(messages, e.MessageTo)
			var s []string
			s = append(s, e.MessageTo)
			err = c.notify.AddNotification(e.UserID, s, 5, 0)
			if err != nil {
				return
			}
		case "group_chat":
			groupNr, err := strconv.Atoi(e.MessageTo)
			if err != nil{
				var errorMessage []dto.WsResponse
				errorMessage[0].Action = "error, invalid group NR"
			}
			messages, groupMembers, err := c.chat.AddGroupMessage(e.UserID, e.MessageTo, e.MessageContent)
			if err != nil {
				return
			}
			for _, member := range groupMembers{
				var s []string
				s = append(s, member.ID)
				err = c.notify.AddNotification(e.UserID, s, 8, groupNr)
				if err != nil {
					return
				}
				c.SendOne(messages, member.ID)
			}
		case "message_read":

		case "left":
			delete(clients, e.Conn)
		}

	}
}

func (c *WsController) getUserList() []string {
	var userList []string
	for _, x := range clients {
		userList = append(userList, x)
	}
	sort.Strings(userList)
	return userList
}

func (c *WsController) broadcastToAll(response dto.WsResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			logger.ErrorLogger.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func (c *WsController) SendOne(response []dto.WsResponse, sendTo string) {
	for client, id := range clients {
		if id == sendTo {
			err := client.WriteJSON(response)
			if err != nil {
				logger.ErrorLogger.Println("websocket err", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}

func (c *WsController) InitRoutes() {
	go c.ListenToWsChannel()
	c.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		//logger.InfoLogger.Println("Client Successfully Connected...")
		var response dto.WsResponse
		response.Action = `Connected to server`
		conn := WebSocketConn{Conn: ws}
		clients[conn] = ""

		err = ws.WriteJSON(response)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		go c.reader(&conn)
	})
}
