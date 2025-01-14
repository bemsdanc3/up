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
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
	query := `INSERT INTO users (login, email, pass) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Login, user.Email, user.Pass)
	if err != nil {
		log.Printf("Error creating user: %v", err)
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
	query := `SELECT id, email, pass FROM users WHERE email = $1`
	var user entities.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Pass)
	if err != nil {
		log.Printf("Error querying user by email (%s): %v", email, err)
		return nil, err
	}
	return &user, nil
}
