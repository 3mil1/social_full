package repo

import (
	"database/sql"
	"social-network/internal/dto"
	"social-network/internal/entity"
)


type PostRepo struct {
	db *sql.DB
	user Repo
}

func NewPostRepo(db *sql.DB, user Repo) *PostRepo {
	return &PostRepo{db, user}
}

func (r PostRepo) AddNewPost(post entity.UserPost) (int, error){
	var parentId *int
	if post.ParentId != 0{
		parentId = &post.ParentId
	}
	query := "INSERT INTO post (user_id, title, content, image, parent_id, privacy) VALUES ($1, $2, $3, $4, $5, $6) returning id"
	row := r.db.QueryRow(query, post.UserId, post.Subject, post.Content, post.Image, parentId, post.Privacy)
	err := row.Scan(&post.Id)
	if err != nil {
		return 0, err
	}
	return post.Id, nil
}

func (r PostRepo) AddNewPostAccess(postId int, oneUserId string) error {
	query := "INSERT INTO post_access (post_id, user_id) VALUES ($1, $2)"
	if _, err := r.db.Exec(query, postId, oneUserId ); err != nil {
		return err
		}
	return nil
}


func (r PostRepo) GetAllUserPosts(loggedInUserId, userId string) ([]dto.PostReply, error){
	var list []dto.PostReply
	rows, err := r.db.Query(`SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, p.updated_at, p.privacy, u.first_name, u.last_name 
	FROM [post] p 
	LEFT JOIN [user] u ON u.id = p.user_id  
	WHERE user_id=? AND parent_id IS NULL
	ORDER BY p.id DESC`, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post dto.PostReply
		err := rows.Scan(&post.Id, &post.UserId, &post.Subject, &post.Content, 
			&post.Image, &post.CreatedAt,&post.UpdatedAt, &post.Privacy, &post.UserFirstName, &post.UserLastName)
		if err != nil {
			return nil, err
		}
		//check if post is private -> is there approved follower connection
		if post.Privacy == 2 && loggedInUserId != userId {
			access, err := r.user.Get2UsersConnectionStatus(loggedInUserId, userId)
			if err != nil{
				return nil, err
			}
			if access == 1{
				list = append(list, post)
			}
		} 
		if post.Privacy == 3 && loggedInUserId != userId {
			//check, that loggedInUser has access
			access, err := r.GetAccessByUserId(loggedInUserId, post.Id)
			if err != nil{
				return nil, err
			}
			if access > 0 {
				list = append(list, post)
			}
		}
		if post.Privacy == 1 || loggedInUserId == post.UserId {
			list = append(list, post)
		}
	}
	return list, nil
}


func (r PostRepo) GetAllPosts(loggedInUserId string) ([]dto.PostReply, error){
	var list []dto.PostReply
	rows, err := r.db.Query(`SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, p.updated_at, p.privacy, u.first_name, u.last_name 
	FROM [post] p LEFT JOIN [user] u ON u.id = p.user_id  
	WHERE parent_id IS NULL
	ORDER BY p.id DESC`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post dto.PostReply
		err := rows.Scan(&post.Id, &post.UserId, &post.Subject, &post.Content, 
			&post.Image, &post.CreatedAt,&post.UpdatedAt, &post.Privacy, &post.UserFirstName, &post.UserLastName)
		if err != nil {
			return nil, err
		}
		//check if post is private -> is there approved follower connection
		if post.Privacy == 2 && loggedInUserId != post.UserId {
			access, err := r.user.Get2UsersConnectionStatus(loggedInUserId, post.UserId)
			if err != nil{
				return nil, err
			}
			if access == 1{
				list = append(list, post)
			}
		} 
		if post.Privacy == 3 && loggedInUserId != post.UserId {
			//check, that loggedInUser has access
			access, err := r.GetAccessByUserId(loggedInUserId, post.Id)
			if err != nil{
				return nil, err
			}
			if access > 0 {
				list = append(list, post)
			}
		}
		if post.Privacy == 1 || loggedInUserId == post.UserId {
			list = append(list, post)
		}
	}
	return list, nil
}


func(r PostRepo) GetAccessByUserId(userId string, postId int) (int, error){
	var result int
	if err := r.db.QueryRow("SELECT id FROM post_access WHERE post_id= ? AND user_id = ?",
		postId, userId).Scan(&result); err != nil {
	   if err == sql.ErrNoRows {
		   return 0, nil
	   }
	   return 0, nil
   }
   return result, nil
} 


func(r PostRepo) GetOnePostsComments(postId int) ([]dto.PostReply, error){
	var result []dto.PostReply
	rows, err := r.db.Query("SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, p.updated_at, p.privacy, u.first_name, u.last_name FROM [post] p LEFT JOIN [user] u ON u.id = p.user_id WHERE parent_id = ?", postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post dto.PostReply
		err := rows.Scan(&post.Id, &post.UserId, &post.Subject, &post.Content, 
			&post.Image, &post.CreatedAt,&post.UpdatedAt, &post.Privacy, &post.UserFirstName, &post.UserLastName)
		if err != nil {
			return nil, err
		}
		result = append(result, post)
	}
	return result, nil
}

func (r PostRepo) GetPostStatusByPostId(postId int) (int, string, error){
	var status int
	var owner string
	 if err := r.db.QueryRow("SELECT privacy, user_id FROM post WHERE id= ?",
	 	postId).Scan(&status, &owner); err != nil {
        if err == sql.ErrNoRows {
            return 0, "",  err
        }
        return 0, "",  err
    }
    return status, owner, nil
}

func (r PostRepo) GetPostAuthorByPostId(postId int) (string, error){
	var user_id string
	 if err := r.db.QueryRow("SELECT user_id FROM post WHERE id= ?",
	 	postId).Scan(&user_id); err != nil {
        if err == sql.ErrNoRows {
            return "", err
        }
        return "", err
    }
    return user_id, nil
}

func (r PostRepo) CheckPostAccessByPostIdAndUserId(postId int, loggedInUser string) (int, error){
	var id int
	 if err := r.db.QueryRow("SELECT id FROM post_access WHERE post_id = ? AND user_id = ?",
	 	postId, loggedInUser).Scan(&id); err != nil {
        if err == sql.ErrNoRows {
            return 0, nil
        }
        return 0, err
    }
    return id, nil
}


func (r PostRepo) GetOnePostByPostId(postId int) (dto.PostReply, error){
	var result dto.PostReply
	if err := r.db.QueryRow("SELECT p.id, p.user_id, p.title, p.content, p.image, p.created_at, p.updated_at, p.privacy, u.first_name, u.last_name FROM [post] p LEFT JOIN [user] u ON u.id = p.user_id WHERE p.id = ?", postId).Scan(&result.Id, &result.UserId, &result.Subject, &result.Content, 
		&result.Image, &result.CreatedAt,&result.UpdatedAt, &result.Privacy, &result.UserFirstName, &result.UserLastName); err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return result, err
	}
	return result, nil

}

func (r PostRepo) GetPostOwner(postId int) (string, error) {
	var postOwner string
	if err := r.db.QueryRow("SELECT user_id FROM post WHERE id=? ", postId).Scan(&postOwner); err != nil {
		if err == sql.ErrNoRows {
			return postOwner, nil
		}
		return postOwner, err
	}
	return postOwner, nil

}