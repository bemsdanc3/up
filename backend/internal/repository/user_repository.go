package repository

import (
	"backend/internal/entities"
	"database/sql"
	"log"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	IsEmailExists(email string) (bool, error)
	IsLoginExists(login string) (bool, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(userID int) (*entities.User, error)
	GetPlaylistsByUserID(userID int) ([]entities.Playlist, error)
	GetSubscriptionsByUserID(userID int) ([]int, error)
	UpdateUserProfile(user *entities.User, userID int) error
	GetUserProfile(userID int) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
	UpdateUsersRole(userID int, role string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
	query := `INSERT INTO users (login, email, pass, pfp) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, user.Login, user.Email, user.Pass, user.PFP).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return err
}

func (r *userRepository) IsEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) IsLoginExists(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`
	err := r.db.QueryRow(query, login).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	query := `SELECT id, email, pass, role FROM users WHERE email = $1`
	var user entities.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Pass, &user.Role)
	if err != nil {
		log.Printf("Error querying user by email (%s): %v", email, err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(userID int) (*entities.User, error) {
	query := `SELECT id, login, email, pfp, role FROM users WHERE id = $1`
	var user entities.User
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Login, &user.Email, &user.PFP, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("error querying user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPlaylistsByUserID(userID int) ([]entities.Playlist, error) {
	query := `SELECT id, title, cover FROM playlists WHERE author_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying playlist by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()
	var playlists []entities.Playlist
	for rows.Next() {
		var playlist entities.Playlist
		if err := rows.Scan(&playlist.ID, &playlist.Title, &playlist.Cover); err != nil {
			log.Printf("error scanning playlist row: %v", err)
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (r *userRepository) GetSubscriptionsByUserID(userID int) ([]int, error) {
	query := `SELECT artist_id FROM users_artists WHERE user_id = $1 AND isLike = true`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying subs by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var subs []int
	for rows.Next() {
		var artistID int
		if err := rows.Scan(&artistID); err != nil {
			log.Printf("error scanning subs row: %v", err)
			return nil, err
		}
		subs = append(subs, artistID)
	}
	return subs, nil
}

func (r *userRepository) UpdateUserProfile(user *entities.User, userID int) error {
	query := `UPDATE users SET login = $1, email = $2, pass = $3, pfp = $4 WHERE id = $5`
	_, err := r.db.Exec(query, user.Login, user.Email, user.Pass, user.PFP, userID)
	if err != nil {
		log.Printf("error updating user on database lvl: %v", err)
		return err
	}
	return err
}

func (r *userRepository) GetUserProfile(userID int) (*entities.User, error) {
	query := `SELECT login, pfp, email, role FROM users WHERE id = $1`
	var user entities.User
	err := r.db.QueryRow(query, userID).Scan(&user.Login, &user.PFP, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("error querying user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]entities.User, error) {
	query := `SELECT id, login, email, role, pfp, creation_date FROM users ORDER BY creation_date`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("error getting users on database lvl: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err = rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.Role,
			&user.PFP,
			&user.CreationDate,
		); err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return nil, err
	}
	return users, err
}

func (r *userRepository) UpdateUsersRole(userID int, role string) error {
	query := `UPDATE users set role = $1 WHERE id = $2`
	result, err := r.db.Exec(query, role, userID)
	if err != nil {
		log.Printf("error updating role on database lvl: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("user not found")
		return err
	}

	return nil
}
