package service

import (
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"
)

type NotificationService struct {
	repo NotificationRepo
	ws   WsController
}

func NewNotificationService(r NotificationRepo, ws WsController) *NotificationService {
	return &NotificationService{
		repo: r,
		ws:   ws,
	}
}

type WsController interface {
	SendOne(response []dto.WsResponse, sendTo string)
}

func (s *NotificationService) AddNotification(actorID string, receiverID []string, notificationType int, objectId int) error {
	var n = entity.Notification{
		ActorID:          actorID,
		NotificationType: notificationType,
		ObjectID:         objectId,
	}

	//check notification object existence ("comment to post" or "new chat in group")
	notificationObjId, err := s.repo.CheckNotifyObj(n)
	if err != nil{
		return err
	}
	if notificationObjId == nil {
		notificationObjId, err = s.repo.AddNotificationOBJ(n)
		if err != nil {
			return err
		}
	}
	n.NotificationObjID = int(*notificationObjId)
	for _, receiver := range receiverID {
		n.ReceiverID = receiver
		_, err := s.repo.AddNotificationReceiver(n)
		if err != nil {
			return err
		}
		notificationsResponse, err := s.GetAllNotifications(receiver)
		if err != nil {
			return err
		}

		if notificationType != 8 && n.NotificationType != 5 {
			s.ws.SendOne(notificationsResponse, n.ReceiverID)	
		} else{
			//for group chat and private chat we send only 1 notification separately
			s.ws.SendOne(notificationsResponse[len(notificationsResponse)-1:], n.ReceiverID)
		}
	}
	return nil
}

func (s *NotificationService) GetAllNotifications(id string) ([]dto.WsResponse, error) {
	notifications, err := s.repo.GetAllNotifications(id)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) DeleteNotification(actorId, receiverId string, notification_type, objectID int) error {
	notifyIds , err := s.repo.FindNotificationId(actorId, receiverId, notification_type, objectID)
	if err != nil {
		return err
	}
	for _, notify_id := range notifyIds{
		err2 := s.repo.DeleteNotification(receiverId, notify_id)
		if err2 != nil{
			logger.ErrorLogger.Println(err)
		}
	}
	return nil
}

func (s *NotificationService) UpdateNotification(loggedInUserId, notify_id, status string) error{
	err := s.repo.UpdateNotification(loggedInUserId, notify_id, status)
	if err != nil{
		return err
	}
	return nil
}

func (s *NotificationService) DeleteGroupChatNotification(loggedInUserId, group_id string) error {
	return s.repo.DeleteGroupChatNotification(loggedInUserId, group_id)
}