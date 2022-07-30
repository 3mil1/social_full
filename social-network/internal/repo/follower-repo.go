package repo

import (
	"database/sql"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"
)

type FRepo struct {
	*sql.DB
}

func NewFRepo(db *sql.DB) *FRepo {
	return &FRepo{db}
}

func (r FRepo) GetAllUsersIFollow(someUserId string) ([]dto.FollowerUserSmall, error) {
	var users []dto.FollowerUserSmall
	rows, err := r.Query("SELECT u.id, u.first_name, u.last_name, u.image, f.status FROM follower f JOIN user u ON u.id = f.target_id  WHERE source_id=? AND status = 1", someUserId)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var oneUser dto.FollowerUserSmall
		if err := rows.Scan(&oneUser.UserId, &oneUser.FirstName, &oneUser.LastName, &oneUser.UserImg, &oneUser.Status); err != nil {
			logger.WarningLogger.Println(err)
			return nil, err
		}
		users = append(users, oneUser)
	}
	return users, nil
}

func (r FRepo) CheckFollowRequest(someUserId, targetId string) (int, error) {
	var status int
	err := r.QueryRow("SELECT status FROM follower WHERE source_id=? AND target_id = ?", someUserId, targetId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return status, nil
}

func (r FRepo) AddNewFollower(follower entity.Follower) error {
	query := "INSERT INTO follower (source_id, target_id, status) VALUES ($1, $2, $3)"
	if _, err := r.Exec(query, follower.SourceId, follower.TargetId, follower.Status); err != nil {
		return err
	}
	return nil
}

func (r FRepo) UpdateFollower(follower entity.Follower) error {
	_, err := r.Exec("UPDATE follower SET status = ? WHERE (source_id = ? AND target_id = ?)", follower.Status, follower.SourceId, follower.TargetId)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (r FRepo) GetUsersFollowsMeAcceptedStatusOnly(targetId string) ([]dto.FollowerUserSmall, error) {
	var usersFollowsMe []dto.FollowerUserSmall
	rows, err := r.Query("SELECT u.id, u.first_name, u.last_name, u.image, f.status FROM follower f JOIN user u ON u.id = f.source_id  WHERE target_id=? AND status = 1", targetId)

	if err != nil {
		return nil, err //"Some problem with reading data from DB"
	}
	for rows.Next() {
		var oneUser dto.FollowerUserSmall
		if err := rows.Scan(&oneUser.UserId, &oneUser.FirstName, &oneUser.LastName, &oneUser.UserImg, &oneUser.Status); err != nil {
			logger.WarningLogger.Println(err)
			return nil, err
		}
		usersFollowsMe = append(usersFollowsMe, oneUser)

	}
	return usersFollowsMe, nil
}
func (r FRepo) GetAllUsersFollowsMe(targetId string) ([]dto.FollowerUserSmall, error) {
	var usersFollowsMe []dto.FollowerUserSmall
	rows, err := r.Query("SELECT u.id, u.first_name, u.last_name, u.image, f.status FROM follower f JOIN user u ON u.id = f.source_id  WHERE target_id=?", targetId)

	if err != nil {
		return nil, err //"Some problem with reading data from DB"
	}
	for rows.Next() {
		var oneUser dto.FollowerUserSmall
		if err := rows.Scan(&oneUser.UserId, &oneUser.FirstName, &oneUser.LastName, &oneUser.UserImg, &oneUser.Status); err != nil {
			logger.WarningLogger.Println(err)
			return nil, err
		}
		usersFollowsMe = append(usersFollowsMe, oneUser)
	}
	return usersFollowsMe, nil
}


func (r FRepo) DeleteFollower(loggedInUserId, userId string) error {
	query := "DELETE FROM follower WHERE source_id = ? and target_id = ?"
	if _, err := r.Exec(query, loggedInUserId, userId); err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return nil
}