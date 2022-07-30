package service

import (
	"encoding/json"
	"social-network/internal/dto"
	"social-network/pkg/logger"
	"strconv"
)

type ChatService struct {
	repo ChatRepo
	groupRepo GroupRepo
}

func NewChatService(r ChatRepo, gr GroupRepo) *ChatService {
	return &ChatService{
		repo: r,
		groupRepo: gr,
	}
}

func (s *ChatService) GetMessages(sender, receiver string, skip, limit int) ([]dto.ChatMessage, error) {
	messages, err := s.repo.GetMessages(sender, receiver, skip, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *ChatService) AddMessage(sender, receiver, messageContent string) ([]dto.WsResponse, error) {
	var message dto.WsResponse
	m, err := s.repo.AddChat(sender, receiver, messageContent)
	if err != nil {
		return nil, err
	}

	message.Action = "message"
	message.ActionType = "private message"

	var myMap map[string]interface{}
	data, _ := json.Marshal(m)
	err = json.Unmarshal(data, &myMap)
	if err != nil {
		return nil, err
	}

	
	message.Data = myMap

	var messages []dto.WsResponse

	messages = append(messages, message)

	return messages, nil
}


func (s *ChatService) AddGroupMessage(sender, groupId, messageContent string) ([]dto.WsResponse, []dto.PrivateProfileResponse, error) {

	var message dto.WsResponse

	m, err := s.repo.AddGroupMessage(sender, groupId, messageContent)
	if err != nil {
		return nil, nil, err
	}

	message.Action = "message"
	message.ActionType = "group message"

	var myMap map[string]interface{}
	data, _ := json.Marshal(m)
	err = json.Unmarshal(data, &myMap)
	if err != nil {
		return nil, nil, err
	}
	groupNr, err := strconv.Atoi(groupId)
	if err != nil{
		return nil, nil, err
	}
	members, err2 := s.groupRepo.GetGroupMemberDetails(groupNr)
	if err2 != nil{
		return nil, nil, err2
	}
	for _, member := range members{
		if member.ID == sender {
			err := s.repo.AddGroupChatReceiver(m.MessageID, member.ID, true)
			if err != nil{
				logger.ErrorLogger.Println(err)
			}
		} else{
			err := s.repo.AddGroupChatReceiver(m.MessageID, member.ID, false)
			if err != nil{
				logger.ErrorLogger.Println(err)
			}
		}
	}

	message.Data = myMap
	var messages []dto.WsResponse
	messages = append(messages, message)

	return messages, members, nil
}

func (s *ChatService) MarkMessageAsSeen(loggedInUserId, target_id string){
	//target is user_id
	if len(target_id) > 20{
		s.repo.MarkPersonalMessageAsSeen(loggedInUserId, target_id)
	}else{
		//target is groupId
		s.repo.MarkGroupMessagesAsSeen(loggedInUserId, target_id)
	}

}