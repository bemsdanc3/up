package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserUseCase interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserDetails(userID int) (*entities.UserDetails, error)
	UpdateUserProfile(user *entities.User, userID int) error
	GetUserProfile(userID int) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUsecse(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) CreateUser(user *entities.User) error {
	emailExists, err := u.repo.IsEmailExists(user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return fmt.Errorf("email is already in use")
	}

	loginExists, err := u.repo.IsLoginExists(user.Login)
	if err != nil {
		return err
	}
	if loginExists {
		return fmt.Errorf("login is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Pass = string(hashedPassword)
	return u.repo.CreateUser(user)
}

func (u *userUseCase) GetUserByEmail(email string) (*entities.User, error) {
	return u.repo.GetUserByEmail(email)
}

func (u *userUseCase) GetUserDetails(userID int) (*entities.UserDetails, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		log.Printf("error getting user by ID: %v", err)
		return nil, err
	}
	subs, err := u.repo.GetSubscriptionsByUserID(userID)
	if err != nil {
		log.Printf("error getting sub by userID: %v", err)
		return nil, err
	}
	playlist, err := u.repo.GetPlaylistsByUserID(userID)
	if err != nil {
		log.Printf("error getting playlist by userID: %v", err)
		return nil, err
	}

	return &entities.UserDetails{
		User:         *user,
		Subscription: subs,
		Playlists:    playlist,
	}, nil
}

func (u *userUseCase) UpdateUserProfile(user *entities.User, userID int) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Pass = string(hashedPassword)
	return u.repo.UpdateUserProfile(user, userID)
}

func (u *userUseCase) GetUserProfile(userID int) (*entities.User, error) {
	return u.repo.GetUserProfile(userID)
}

func (u *userUseCase) GetAllUsers() ([]entities.User, error) {
	return u.repo.GetAllUsers()
}
