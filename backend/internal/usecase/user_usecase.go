package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserUsecase interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecse(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(user *entities.User) error {
	emailExists, err := u.repo.IsEmailExists(user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		log.Printf("Email is already in use")
	}

	loginExists, err := u.repo.IsLoginExists(user.Login)
	if err != nil {
		return err
	}
	if loginExists {
		log.Printf("Login is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Pass = string(hashedPassword)
	return u.repo.CreateUser(user)
}

func (u *userUsecase) GetUserByEmail(email string) (*entities.User, error) {
	return u.repo.GetUserByEmail(email)
}
