package test

import (
	"errors"
	mocks "monorepo-ecommerce/micro-services/user/mocks/mock_micro-services/user/repository"
	"monorepo-ecommerce/micro-services/user/models"
	"monorepo-ecommerce/micro-services/user/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userService := service.NewUserService(mockRepo)

	email := "test@example.com"
	phone := "1234567890"
	password := "securepassword"

	user := &models.User{
		Id:       1,
		Email:    email,
		Phone:    phone,
		Password: password,
	}

	t.Run("should success", func(t *testing.T) {
		mockRepo.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

		user, err := userService.RegisterUser(email, phone, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, phone, user.Phone)
	})

	t.Run("should error validation", func(t *testing.T) {
		_, err := userService.RegisterUser("", "", "password")
		assert.Error(t, err)
		assert.EqualError(t, err, "email or phone is required")
	})
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userService := service.NewUserService(mockRepo)

	email := "test@example.com"
	phone := "1234567890"
	password := "securepassword"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Id:       1,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	t.Run("should success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmailOrPhone(email, phone, password).Return(user, nil)

		user, err := userService.LoginUser(email, phone, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, phone, user.Phone)
	})

	t.Run("should error when user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmailOrPhone(email, phone, password).Return(nil, errors.New("user not found"))

		user, err := userService.LoginUser(email, phone, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "user not found")
	})
}
