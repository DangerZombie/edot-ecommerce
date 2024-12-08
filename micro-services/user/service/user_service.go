package service

import (
	"errors"
	"fmt"
	"log"
	"monorepo-ecommerce/micro-services/user/models"
	"monorepo-ecommerce/micro-services/user/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(email string, phone string, password string) (*models.User, error)
	LoginUser(email string, phone string, password string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(email string, phone string, password string) (user *models.User, err error) {
	if email == "" || phone == "" {
		return nil, errors.New("email or phone is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v\n", err)
		return nil, err
	}

	data := models.User{
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	user, err = s.repo.CreateUser(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) LoginUser(email string, phone string, password string) (user *models.User, err error) {
	user, err = s.repo.GetUserByEmailOrPhone(email, phone, password)
	fmt.Println(user)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
