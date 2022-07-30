package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type ChatRepo struct {
	*sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db}
}

func (r *ChatRepo) CountMessages(sender, receiver string) (int, error) {
	row := r.QueryRow(`SELECT count()
		FROM chat c
        LEFT JOIN chat_message m on c.id = m.chat_id
		WHERE (c.user1_id=$1 AND c.user2_id=$2)
		OR (c.user1_id=$2 AND c.user2_id=$1) `, sender, receiver)

	var nOfMessages int

	err := row.Scan(&nOfMessages)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	return nOfMessages, nil
}

func (r *ChatRepo) GetMessages(sender, receiver string, skip, limit int) ([]dto.ChatMessage, error) {
	rows, err := r.Query(`SELECT m.user_id, m.content, m.created_at
		FROM chat c
        LEFT JOIN chat_message m on c.id = m.chat_id
		WHERE (c.user1_id=$1 AND c.user2_id=$2)
		OR (c.user1_id=$2 AND c.user2_id=$1)
		ORDER BY m.id DESC
		LIMIT $3, $4`, sender, receiver, skip, limit)

	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()

	var messages []dto.ChatMessage
	for rows.Next() {
		var m dto.ChatMessage
		err = rows.Scan(&m.From, &m.Content, &m.CreatedAt)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		m.CreatedAt = FormateDate(m.CreatedAt)
		messages = append([]dto.ChatMessage{m}, messages...)
	}
	return messages, nil
}

func (r *ChatRepo) AddChat(sender, receiver, messageContent string) (*entity.ChatMessage, error) {
	var chatMessage entity.ChatMessage

	row := r.QueryRow(`SELECT id
              FROM chat
              WHERE (user1_id=$1
                and user2_id=$2) OR (user1_id=$2 AND user2_id=$1)
`, sender, receiver)

	err := row.Scan(&chatMessage.ChatID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
		}
	}

	if chatMessage.ChatID > 0 {
		err = r.AddMessage(&chatMessage, sender, messageContent, chatMessage.ChatID, receiver)
		if err != nil {
			return nil, err
		}

	} else {
		//create chat if not exist
		var result sql.Result
		if result, err = r.Exec(`INSERT into chat (user1_id, user2_id) VALUES ($1, $2)`, sender, receiver); err != nil {
			logger.WarningLogger.Println("DB error: ", err)
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		err = r.AddMessage(&chatMessage, sender, messageContent, int(id), receiver)
		if err != nil {
			return nil, err
		}

	}
	return &chatMessage, nil
}

func (r *ChatRepo) AddMessage(chatMessage *entity.ChatMessage, sender, messageContent string, chatID int, receiver string) error {

	query := fmt.Sprintf("INSERT into chat_message (%s) VALUES ($1,$2,$3) RETURNING *", "chat_id, user_id, content")
	row := r.QueryRow(query, chatID, sender, messageContent)
	err := row.Scan(&chatMessage.MessageID, &chatMessage.ChatID, &chatMessage.From, &chatMessage.Content, &chatMessage.CreatedAt)
	if err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return err
		}
		logger.ErrorLogger.Println(err)
		return err
	}

	if _, err = r.Exec(`INSERT INTO message_status (message_id, user_id)
		VALUES ($1, $2)`, chatMessage.MessageID, receiver); err != nil {
		logger.WarningLogger.Println("DB error: ", err)
		return err
	}

	return nil
}



func (r *ChatRepo) AddGroupMessage(sender, groupId, messageContent string) (*entity.GroupChatMessage, error){
	var message  entity.GroupChatMessage
	query := fmt.Sprintf("INSERT into group_message (%s) VALUES ($1,$2,$3) RETURNING *", "group_id, user_id, content")
	row := r.QueryRow(query, groupId, sender, messageContent)
	err := row.Scan(&message.MessageID, &message.GroupId, &message.FromId, &message.Content, &message.CreatedAt)
	if err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	message.CreatedAt = FormateDate(message.CreatedAt)

	query2 := "SELECT first_name, last_name FROM user WHERE id = ?"
	row2 := r.QueryRow(query2, sender)
	err2 := row2.Scan(&message.FirstName, &message.LatsName)
	if err2 != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return &message, nil
}

func (r *ChatRepo) AddGroupChatReceiver(messageId int, receiverId string, read bool) error {
	_, err := r.Exec(`INSERT INTO group_message_status (message_id, user_id, read)
	VALUES ($1, $2, $3)`, messageId, receiverId, read)
	if err != nil {
		logger.WarningLogger.Println("DB error: ", err)
		return err
	}
	return nil

}

func FormateDate(date string) string{
	splited := strings.Split(date, "T")
	time := splited[1][:len(splited[1])-4]
	dates := strings.Split(splited[0], "-")
	new_date := dates[2] + "-" +  dates[1] + "-" + dates[0]
	return time + " "+ new_date
}

func (r *ChatRepo) MarkPersonalMessageAsSeen(loggedInUserId, target_id string) {
	row := r.QueryRow(`SELECT id FROM chat WHERE (user1_id=$1 AND user2_id=$2)
	OR (user1_id=$2 AND user2_id=$1)`, loggedInUserId, target_id)
	var id int
	err := row.Scan(&id)
	if err != nil{
		logger.ErrorLogger.Println(err)
	}

	var message_ids []int
	rows, err := r.Query("SELECT id FROM chat_message WHERE chat_id = ?", id)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	for rows.Next() {
		var message_id int
		err := rows.Scan(&message_id)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		message_ids = append(message_ids, message_id)
	}

	for _, id := range message_ids{
		_, err2 := r.Exec(`DELETE FROM message_status WHERE user_id = ? AND message_id = ?`, loggedInUserId, id)
		if err2 != nil {
			logger.ErrorLogger.Println(err)
		}	
	}
}

func (r *ChatRepo) MarkGroupMessagesAsSeen(loggedInUserId, target_id string) {
	rows, err := r.Query(`SELECT id FROM group_message WHERE group_id = ?`, target_id)
	if err != nil{
		logger.ErrorLogger.Println(err)
	}

	var message_ids []int
	for rows.Next() {
		var message_id int
		err := rows.Scan(&message_id)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		message_ids = append(message_ids, message_id)
	}

	for _, id := range message_ids{
		_, err2 := r.Exec(`UPDATE group_message_status SET read = true WHERE user_id = ? AND message_id = ?`, loggedInUserId, id)
		if err2 != nil {
			logger.ErrorLogger.Println(err)
		}
	}
}