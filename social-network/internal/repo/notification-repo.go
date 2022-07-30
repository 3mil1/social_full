package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"
)

type NotificationRepo struct {
	*sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db}
}

func (r *NotificationRepo) AddNotificationReceiver(n entity.Notification) (*entity.Notification, error) {
	query := fmt.Sprintf("INSERT INTO notification (%s) VALUES ($1, $2, $3)", "receiver_id, notification_id, seen")
	if _, err := r.Exec(query, n.ReceiverID, n.NotificationObjID, 0); err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepo) AddNotificationOBJ(n entity.Notification) (*int, error) {
	query := fmt.Sprintf("INSERT INTO notification_obj (%s) VALUES ($1, $2, $3) returning id", "notification_type, object_id, actor_id")
	row := r.QueryRow(query, n.NotificationType, n.ObjectID, n.ActorID)
	err := row.Scan(&n.ObjectID)
	if err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return &n.ObjectID, nil
}

func (r *NotificationRepo) GetAllNotifications(id string) ([]dto.WsResponse, error) {
	query := `SELECT n.id, n.seen, nt.type, nOBJ.actor_id, nOBJ.object_id, u.first_name, u.last_name
	FROM notification_obj nOBJ
	JOIN notification n on nOBJ.id = n.notification_id
	JOIN notification_type nt on nOBJ.notification_type = nt.id
	JOIN user u ON u.id = nOBJ.actor_id
	WHERE n.receiver_id=$1
	`

	rows, err := r.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []dto.WsResponse
	for rows.Next() {
		n := dto.WsResponse{}
		n.Action = "notification"
		var notify_id int
		var actorId string
		var firstName string
		var lastName string
		var objectID int
		var seen int
		err = rows.Scan(&notify_id, &seen, &n.ActionType, &actorId, &objectID,  &firstName, &lastName)
		m := map[string]interface{}{}
		m["notif_id"] = notify_id
		m["actor_id"] = actorId
		m["first_name"] = firstName
		m["last_name"] = lastName
		m["seen"] = seen
		
		if n.ActionType == "group invitation" || n.ActionType == "new group member request"|| 
		n.ActionType == "group access opened" || n.ActionType == "new message in group chat"{
			var groupTitle string
			err:= r.QueryRow("SELECT title FROM groups where id=?", objectID).Scan(&groupTitle)
			if err != nil{
				logger.ErrorLogger.Println(err)
				return nil, err
			}
			m["group_id"] = objectID
			m["group_name"] = groupTitle
		}
		if n.ActionType == "new event" {
			var groupTitle string
			var groupId int
			var eventName string
			err:= r.QueryRow("SELECT e.group_id, e.title, g.title FROM group_event e JOIN groups g ON g.id = e.group_id WHERE e.id=?", objectID).Scan(&groupId, &eventName, &groupTitle)
			if err != nil{
				logger.ErrorLogger.Println(err)
				return nil, err
			}
			m["group_id"] = groupId
			m["group_name"] = groupTitle
			m["event_id"] = objectID
			m["event_name"] = eventName
		}
		if n.ActionType == "new comment to post" {
			var postTitle string
			err:= r.QueryRow("SELECT title FROM post where id=?", objectID).Scan(&postTitle)
			if err != nil{
				return nil, err
			}
			m["post_id"] = objectID
			m["post_name"] = postTitle
		}

		n.Data = m
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepo) DeleteNotification(receiverId string, notifyId int) error{
	sqlStatement := "DELETE FROM notification WHERE receiver_id = $1 AND notification_id = $2"
	_, err := r.Exec(sqlStatement, receiverId, notifyId)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (r *NotificationRepo) FindNotificationId(actorId, receiverId string, notification_type, objectID int) ([]int, error){
	var notifyId []int
	rows, err := r.Query("SELECT id FROM notification_obj WHERE notification_type = ? AND object_id = ? AND actor_id = ?", notification_type, objectID, actorId)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	for rows.Next(){
		var oneId int
		err_scan := rows.Scan(&oneId)
		if err_scan != nil{
			logger.ErrorLogger.Println(err)
		}
		notifyId = append(notifyId, oneId)
	}
	return notifyId, nil
}

func (r *NotificationRepo) UpdateNotification(loggedInUserId, notify_id, status string) error{
	query := "UPDATE notification SET seen = ? WHERE receiver_id = ? AND id = ? "
	if _, err := r.Exec(query, status, loggedInUserId, notify_id); err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepo) DeleteGroupChatNotification(loggedInUserId, group_id string) error {
	query := "SELECT id FROM notification_obj WHERE notification_type = 8 AND object_id = ?"
	rows, err := r.Query(query, group_id)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	for _, id := range ids{
		err := r.DeleteNotification(loggedInUserId, id)
		if err != nil{
			logger.ErrorLogger.Println(err)
		}
	}
	return nil
}

func (r *NotificationRepo) CheckNotifyObj(n entity.Notification) (*int, error) {
	query := `SELECT id FROM notification_obj WHERE notification_type = ? AND object_id = ? AND actor_id = ?`
	row := r.QueryRow(query, n.NotificationType, n.ObjectID, n.ActorID)
	err := row.Scan(&n.NotificationObjID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return &n.NotificationObjID, nil
}