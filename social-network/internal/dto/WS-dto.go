package dto

type WsResponse struct {
	Action     string                 `json:"action"`      // notification || message
	ActionType string                 `json:"action_type"` // friend request || private message
	Data       map[string]interface{} `json:"data"`        //
}
