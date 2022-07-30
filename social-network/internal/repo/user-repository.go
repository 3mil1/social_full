package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/internal/dto"
	"social-network/internal/entity"
	"social-network/pkg/logger"

	"github.com/mattn/go-sqlite3"
)

var (
	userColumns = "id, email, password_hash, first_name, last_name, birthday, image, nickname, about, is_private"
)

type Repo struct {
	*sql.DB
}

func NewUserRepo(db *sql.DB) *Repo {
	return &Repo{db}
}

func (r Repo) AddUser(entity *entity.User) (*entity.User, error) {
	query := fmt.Sprintf("INSERT INTO user (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", userColumns)
	if _, err := r.Exec(query, entity.ID, entity.Email, entity.PwHash, entity.FirstName, entity.LastName, entity.DateOfBirth, entity.UserImg, entity.Nickname, entity.AboutMe, entity.IsPrivate); err != nil {
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return entity, nil
}

func (r Repo) GetUserByID(id string) (*entity.User, error) {
	row := r.QueryRow("SELECT * FROM user WHERE id=?", id)

	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.PwHash, &user.FirstName, &user.LastName, &user.DateOfBirth,
		&user.UserImg, &user.Nickname, &user.AboutMe,
		&user.CreatedAt, &user.UpdatedAt, &user.IsPrivate)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	user.DateOfBirth = user.DateOfBirth[:10]
	return &user, nil
}

func (r Repo) GetUserByEmail(email string) (*entity.User, error) {
	row := r.QueryRow("SELECT * FROM user WHERE email=?", email)

	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.PwHash, &user.FirstName, &user.LastName, &user.DateOfBirth,
		&user.UserImg, &user.Nickname, &user.AboutMe,
		&user.CreatedAt, &user.UpdatedAt, &user.IsPrivate)

	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r Repo) AddRefreshToken(refreshToken entity.RefreshTokenDB) error {
	_, err := r.Exec("INSERT INTO sessions VALUES(?,?,?,?);", refreshToken.RefreshToken, refreshToken.UserID, refreshToken.Device, refreshToken.IP)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (r Repo) DeleteRefreshToken(refreshToken entity.RefreshTokenDB) (*entity.RefreshTokenDB, error) {
	row := r.QueryRow("SELECT * FROM sessions WHERE refresh_token=?", refreshToken.RefreshToken)
	var rt entity.RefreshTokenDB

	err := row.Scan(&rt.RefreshToken, &rt.UserID, &rt.Device, &rt.IP)
	if err != nil {
		logger.ErrorLogger.Println(err)

		if err.Error() == "sql: no rows in result set" {
			// deleting all user sessions
			stmt, err := r.Prepare("DELETE FROM sessions WHERE user_id=?")
			if err != nil {
				logger.ErrorLogger.Println(err)
				return nil, err
			}

			_, err = stmt.Exec(refreshToken.UserID)
			if err != nil {
				logger.ErrorLogger.Println(err)
				return nil, err
			}
			logger.WarningLogger.Println("send email: seems someone stoles your identity, we are logging you out from all devices")
			return nil, errors.New("seems someone stoles your identity, we are logging you out from all devices")
		}

		return nil, err
	}

	stmt, err := r.Prepare("DELETE FROM sessions WHERE refresh_token=?")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	res, err := stmt.Exec(refreshToken.RefreshToken)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return &rt, nil
}

func (r Repo) DeleteSession(userID, ip, userAgent string) error {
	stmt, err := r.Prepare("DELETE FROM sessions WHERE user_id=? AND device=? AND ip=?")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	res, err := stmt.Exec(userID, userAgent, ip)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (r Repo) UpdateUser(user dto.UserUpdate, loggedInUser string) error {
	_, err := r.Exec("UPDATE user SET nickname = ?, about = ?, image = ?, is_private = ?, updated_at = (datetime('now','localtime')) WHERE id = ?", user.Nickname, user.AboutMe, user.UserImg, user.IsPrivate, loggedInUser)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	} else {
		return nil
	}
}

func (r Repo) GetUserStatusByID(requestedUserId string) (bool, error) {
	var status bool
	if err := r.QueryRow("SELECT is_private FROM user WHERE id= ?",
		requestedUserId).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no such user in DB")
		}
		return false, fmt.Errorf("some error in DB GetUserStatusByID")
	}
	return status, nil
}

func (r Repo) Get2UsersConnectionStatus(loggedInUserId, requestedUserId string) (int, error) {
	var connection int
	if err := r.QueryRow("SELECT status FROM follower WHERE (source_id= ? AND target_id = ?)",
		loggedInUserId, requestedUserId).Scan(&connection); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		logger.ErrorLogger.Println(err)
		return 0, fmt.Errorf("some error in DB GetUserStatusByID")
	}
	return connection, nil
}

func (r Repo) GetAllUsers() ([]dto.PrivateProfileResponse, error) {
	var result []dto.PrivateProfileResponse
	rows, err := r.Query("SELECT id, first_name, last_name, image FROM user")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user dto.PrivateProfileResponse
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserImg)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}
