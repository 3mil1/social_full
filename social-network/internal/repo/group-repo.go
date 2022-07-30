package repo

import (
	"database/sql"
	"fmt"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"
	"strings"
)

type GroupRepo struct {
	*sql.DB
}

func NewGroupRepo(db *sql.DB) *GroupRepo {
	return &GroupRepo{db}
}

func (r *GroupRepo) AddNewGroup(group entity.Group) (int, error) {
	query := "INSERT INTO [groups] (created_by, title, description) VALUES ($1, $2, $3)  returning id"
	row := r.QueryRow(query, group.Creator, group.Title, group.Description)
	err := row.Scan(&group.Id)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return 0, err
	}
	err2 := r.AddGroupMember(group.Id, group.Creator, "admin", 1)
	if err2 != nil {
		logger.ErrorLogger.Println(err)
		return 0, err2
	}
	return group.Id, nil
}

func (r *GroupRepo) GetAllGroups() ([]dto.GroupShortInfo, error) {
	var list []dto.GroupShortInfo
	rows, err := r.Query("SELECT g.id, g.title, g.description FROM [groups] g ")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var group dto.GroupShortInfo
		err := rows.Scan(&group.Id, &group.Title, &group.Description)
		if err != nil {
			return nil, err
		}
		members, err := r.GetGroupMemberNumber(group.Id)
		if err != nil {
			return nil, err
		}
		group.Members = members
		list = append(list, group)
	}
	return list, nil
}

func (r *GroupRepo) GetAllMyCreatedGroups(loggedInUser string) ([]dto.GroupShortInfo, error) {
	var list []dto.GroupShortInfo
	rows, err := r.Query("SELECT g.id, g.title, g.description FROM [groups] g WHERE created_by = ?", loggedInUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var group dto.GroupShortInfo
		err := rows.Scan(&group.Id, &group.Title, &group.Description)
		if err != nil {
			return nil, err
		}
		members, err := r.GetGroupMemberNumber(group.Id)
		if err != nil {
			return nil, err
		}
		group.Members = members
		list = append(list, group)
	}
	
	return list, nil
}

func (r *GroupRepo) GetAllMyJoinedGroups(loggedInUser string) ([]dto.GroupShortInfo, error) {
	var list []dto.GroupShortInfo
	rows, err := r.Query("SELECT m.group_id, g.title, g.description FROM [group_member] m LEFT JOIN [groups] g ON g.id = m.group_id WHERE (m.user_id = ? AND status = 1 AND role != 'admin' )", loggedInUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var group dto.GroupShortInfo
		err := rows.Scan(&group.Id, &group.Title, &group.Description)
		if err != nil {
			return nil, err
		}
		members, err := r.GetGroupMemberNumber(group.Id)
		if err != nil {
			return nil, err
		}
		group.Members = members
		list = append(list, group)
	}
	return list, nil
}

func (r *GroupRepo) AddGroupMember(groupId int, userId, role string, status int) error {
	query := "INSERT INTO group_member (group_id, user_id, role, status) VALUES ($1, $2, $3, $4) "
	if _, err := r.Exec(query, groupId, userId, role, status); err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) UpdateGroupMember(groupId int, userId, role string, status int) error {
	query := "UPDATE group_member SET status = ? WHERE group_id = ? AND user_id = ? "
	if _, err := r.Exec(query, status, groupId, userId); err != nil {
		return err
	}
	return nil
}


func (r *GroupRepo) GetPendingJoinRequests(groupId int) ([]dto.GroupAccessRequestUser, error) {
	var list []dto.GroupAccessRequestUser
	rows, err := r.Query("SELECT m.group_id, m.status, m.user_id, m.role, u.first_name, u.last_name FROM [group_member] m LEFT JOIN [user] u ON u.id = m.user_id WHERE m.group_id = ? AND status = 3", groupId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user dto.GroupAccessRequestUser
		err := rows.Scan(&user.GroupId, &user.Status, &user.UserId, &user.Role, &user.UserFirstName, &user.UserLastName)
		if err != nil {
			return nil, err
		}
		list = append(list, user)
	}
	return list, nil

}

func (r *GroupRepo) GetGroupMemberNumber(groupId int) (int, error) {
	var members int
	err := r.QueryRow("SELECT COUNT(*) FROM group_member WHERE group_id = ? AND status = 1", groupId).Scan(&members)
	if err != nil {
		return 0, err
	}
	return members, nil
}

func (r *GroupRepo) GetGroupMemberDetails(groupId int) ([]dto.PrivateProfileResponse, error) {
	var list []dto.PrivateProfileResponse
	rows, err := r.Query("SELECT m.user_id, u.first_name, u.last_name, u.image FROM [group_member] m LEFT JOIN [user] u ON u.id = m.user_id WHERE m.group_id = ? AND status = 1", groupId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user dto.PrivateProfileResponse
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserImg)
		if err != nil {
			return nil, err
		}
		list = append(list, user)
	}
	return list, nil

}

func (r *GroupRepo) GetOneGroupInfo(groupId int) (dto.GroupDetailedInfo, error) {
	var group dto.GroupDetailedInfo
	members, err := r.GetGroupMemberDetails(groupId)
	if err != nil {
		return group, err
	}
	group.Members = members

	err = r.QueryRow("SELECT g.id, g.created_by, g.title, g.description, u.first_name, u.last_name FROM [groups] g LEFT JOIN [user] u ON u.id = g.created_by WHERE g.id = ?",
		groupId).Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description, &group.CreatorFirstName, &group.CreatorLastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return group, fmt.Errorf("no group with such ID")
		}
		return group, err
	}
	return group, nil
}

func (r *GroupRepo) AddNewGroupPost(post entity.GroupPost) error {
	var parentId *int
	if post.ParentId != 0 {
		parentId = &post.ParentId
	}
	query := "INSERT INTO group_post (group_id, user_id, title, content, image, parent_id) VALUES ($1, $2, $3, $4, $5, $6) "
	if _, err := r.Exec(query, post.GroupId, post.CreatorId, post.Title, post.Content, post.Image, parentId); err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) GetOneGroupAllPosts(groupId int) ([]dto.GroupPostReply, error) {
	var posts []dto.GroupPostReply
	rows, err := r.Query("SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, u.first_name, u.last_name FROM [group_post] p LEFT JOIN [user] u ON u.id = p.user_id WHERE p.group_id = ? AND parent_id IS NULL",
		groupId)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, fmt.Errorf("no group with such ID")
		}
		return posts, err
	}
	for rows.Next() {
		var post dto.GroupPostReply
		err := rows.Scan(&post.PostId, &post.UserId, &post.Subject, &post.Content, &post.Image, &post.CreatedAt, &post.UserFirstName, &post.UserLastName)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = strings.Replace(post.CreatedAt, "T", " ", 1)[:len(post.CreatedAt)-4]
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *GroupRepo) CheckUserIsGroupMember(loggedInUser string, groupId int) (bool, error) {
	var member int
	err := r.QueryRow("SELECT id FROM group_member WHERE group_id = ? AND user_id = ? AND status = 1", groupId, loggedInUser).Scan(&member)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if member > 0 {
		return true, nil
	}
	return false, nil
}

func (r *GroupRepo) GetOnePost(groupId, postId int) (dto.GroupPostReply, error) {
	var post dto.GroupPostReply
	err := r.QueryRow("SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, u.first_name, u.last_name FROM [group_post] p LEFT JOIN [user] u ON u.id = p.user_id WHERE p.group_id = ? AND p.id = ?",
		groupId, postId).Scan(&post.PostId, &post.UserId, &post.Subject, &post.Content, &post.Image, &post.CreatedAt, &post.UserFirstName, &post.UserLastName)
	if err != nil {
		return post, err
	}
	post.CreatedAt = strings.Replace(post.CreatedAt, "T", " ", 1)[:len(post.CreatedAt)-4]
	return post, nil
}

func (r *GroupRepo) GetPostComments(groupId, postId int) ([]dto.GroupPostReply, error) {
	var posts []dto.GroupPostReply
	rows, err := r.Query("SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, u.first_name, u.last_name FROM [group_post] p LEFT JOIN [user] u ON u.id = p.user_id WHERE p.group_id = ? AND p.parent_id = ?",
		groupId, postId)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post dto.GroupPostReply
		err := rows.Scan(&post.PostId, &post.UserId, &post.Subject, &post.Content, &post.Image, &post.CreatedAt, &post.UserFirstName, &post.UserLastName)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = strings.Replace(post.CreatedAt, "T", " ", 1)[:len(post.CreatedAt)-4]

		posts = append(posts, post)
	}
	return posts, nil
}

func (r *GroupRepo) ReplyToGroupInvitation(loggedInUser string, status, groupId int) error {
	query := "UPDATE group_member SET status = ? WHERE group_id = ? AND user_id = ?"
	if _, err := r.Exec(query, status, groupId, loggedInUser); err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) CheckUserStatusInGroup(loggedInUser string, groupId int) (int, error) {
	var status int
	err := r.QueryRow("SELECT status FROM group_member WHERE group_id = ? AND user_id = ?", groupId, loggedInUser).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return status, nil
}

func (r *GroupRepo) CheckUserRoleInGroup(loggedInUser string, groupId int) (bool, error) {
	var role string
	err := r.QueryRow("SELECT role FROM group_member WHERE group_id = ? AND user_id = ?", groupId, loggedInUser).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if role == "admin" {
		return true, nil
	}
	return false, nil
}

func (r *GroupRepo) DeleteInvitation(loggedInUserId string, groupId int) error {
	sqlStatement := "DELETE FROM group_member WHERE group_id = $1 AND user_id = $2"
	_, err := r.Exec(sqlStatement, groupId, loggedInUserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) CreateEvent(event entity.GroupEventEntity) (int, error) {
	query := "INSERT INTO group_event (group_id, user_id, title, description, event_date) VALUES ($1, $2, $3, $4, $5) returning id"
	row := r.QueryRow(query, event.GroupId, event.UserId, event.Title, event.Description, event.EventDate)
	err := row.Scan(&event.Id)
	if err != nil {
		return 0, err
	}

	//add event creator to event_participants
	err = r.AddEventParticipant(event.UserId, event.Id, event.GoingStatus)
	if err != nil {
		return 0, err
	}
	return event.Id, nil
}

func (r *GroupRepo) AddEventParticipant(userId string, eventId, option int) error {
	//OPTION:
	// 	1) going
	//  2) not going
	//  3) interested
	query := "INSERT INTO event_participants (event_id, user_id, option) VALUES ($1, $2, $3)"
	if _, err := r.Exec(query, eventId, userId, option); err != nil {
		return err
	}
	return nil
}

func (r *GroupRepo) UpdateEventParticipant(userId string, eventId, option int) error {
	//OPTION:
	// 	1) going
	//  2) not going
	//  3) interested
	query := "UPDATE event_participants SET option = ? WHERE event_id = ? and user_id = ?"
	if _, err := r.Exec(query, option, eventId, userId); err != nil {
		return err
	}
	return nil
}



func (r *GroupRepo) GetGroupIdFromEvent(eventId int) (int, error) {
	var groupId int
	err := r.QueryRow("SELECT group_id FROM group_event WHERE id = ?", eventId).Scan(&groupId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return groupId, nil
}

func (r *GroupRepo) CheckUserGoingInEvent(eventId int, userId string) (int, error) {
	var going int
	err := r.QueryRow("SELECT option FROM event_participants WHERE event_id = ? AND user_id = ?", eventId, userId).Scan(&going)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return going, nil
}

func (r *GroupRepo) GetAllGroupEvents(groupId int) ([]dto.GroupEventReply, error) {
	var list []dto.GroupEventReply
	rows, err := r.Query("SELECT e.id, e.group_id, e.user_id, e.title, e.description, e.event_date, e.created_at, u.first_name, u.last_name FROM [group_event] e LEFT JOIN [user] u ON u.id = e.user_id WHERE (event_date > date('now','localtime') AND group_id = ?)", groupId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event dto.GroupEventReply
		err := rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.Day,
			&event.CreatedAt, &event.CreatorFirstName, &event.CreatorLastName)
		if err != nil {
			return nil, err
		}
		event.Time = event.Day[11:16]
		event.Day = event.Day[:10]
		event.CreatedAt = event.CreatedAt[:10]
		list = append(list, event)
	}
	return list, nil
}

func (r *GroupRepo) GetMyCreatedEvents(loggedInUserId string) ([]dto.GroupEventReply, error) {
	var list []dto.GroupEventReply
	rows, err := r.Query("SELECT e.id, e.group_id, e.user_id, e.title, e.description, e.event_date, e.created_at, u.first_name, u.last_name FROM [group_event] e LEFT JOIN [user] u ON u.id = e.user_id WHERE user_id =? AND event_date > date('now','localtime')", loggedInUserId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event dto.GroupEventReply
		err := rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.Day,
			&event.CreatedAt, &event.CreatorFirstName, &event.CreatorLastName)
		if err != nil {
			return nil, err
		}
		event.Time = event.Day[11:16]
		event.Day = event.Day[:10]
		event.CreatedAt = event.CreatedAt[:10]
		list = append(list, event)
	}
	return list, nil
}

func (r *GroupRepo) GetMyJoinedEvents(loggedInUserId string) ([]dto.GroupEventReply, error) {
	var list []dto.GroupEventReply
	rows, err := r.Query("SELECT p.event_id, e.group_id, e.user_id, e.title, e.description, e.event_date, e.created_at, u.first_name, u.last_name FROM [event_participants] p LEFT JOIN [group_event] e ON p.event_id = e.id LEFT JOIN [user] u ON u.id = e.user_id WHERE p.user_id =? AND event_date > date('now','localtime') AND option = 1", loggedInUserId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event dto.GroupEventReply
		err := rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.Day,
			&event.CreatedAt, &event.CreatorFirstName, &event.CreatorLastName)
		if err != nil {
			return nil, err
		}
		event.Time = event.Day[11:16]
		event.Day = event.Day[:10]
		event.CreatedAt = event.CreatedAt[:10]
		list = append(list, event)
	}
	return list, nil
}

func (r *GroupRepo) GetOneGroupAllMessages(loggedInUser string, groupId, skip, limit int) ([]dto.GroupMessageReply, error) {
	var list []dto.GroupMessageReply

	query := `SELECT m.id, m.group_id, m.user_id, m.content, m.created_at, u.first_name, u.last_name, s.read
		FROM group_message m 
		JOIN user u ON u.id = m.user_id
		JOIN group_message_status s ON s.message_id = m.id AND s.user_id = m.user_id
		WHERE m.group_id = ?
		ORDER BY m.id DESC
		LIMIT $2, $3`
	rows, err := r.Query(query, groupId, skip, limit)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message dto.GroupMessageReply
		err := rows.Scan( &message.MessageId, &message.GroupId, &message.FromId, &message.Content,
			&message.CreatedAt, &message.FirstName, &message.LastName, &message.Seen)
		if err != nil {
			logger.ErrorLogger.Println(err)
			continue
		}
		message.CreatedAt = strings.Replace(message.CreatedAt, "T", " ", 1)[:len(message.CreatedAt)-1]
		list = append([]dto.GroupMessageReply{message}, list...)
	}

	return list, nil
}

func (r *GroupRepo) GetGroupAdmin(groupId int) (string, error) {
	var adminId string
	err := r.QueryRow("SELECT created_by FROM groups WHERE id = ?", groupId).Scan(&adminId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return adminId, nil
}

func (r *GroupRepo) GetAllGroupMembersExceptMe(loggedInUserId string, groupId int) ([]string, error) {
	var list []string
	rows, err := r.Query("SELECT user_id FROM group_member WHERE group_id = ? AND user_id != ? AND status = 1", groupId, loggedInUserId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var userId string
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		list = append(list, userId)
	}

	return list, nil
}
