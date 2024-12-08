package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"monorepo-ecommerce/micro-services/user/handler"
	mocks "monorepo-ecommerce/micro-services/user/mocks/mock_micro-services/user/service"
	"monorepo-ecommerce/micro-services/user/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserService := mocks.NewMockUserService(ctrl)
	h := handler.NewUserHandler(mockUserService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.UserRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockUser := models.User{
			Id:    1,
			Email: reqBody.Email,
			Phone: reqBody.Phone,
		}

		mockUserService.EXPECT().
			RegisterUser(reqBody.Email, reqBody.Phone, reqBody.Password).
			Return(&mockUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RegisterUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": 123,
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RegisterUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should internal server error when something when wrong in RegisterUser", func(t *testing.T) {
		reqBody := handler.UserRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockUserService.EXPECT().
			RegisterUser(reqBody.Email, reqBody.Phone, reqBody.Password).
			Return(nil, errors.New("internal server error"))

		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RegisterUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserService := mocks.NewMockUserService(ctrl)
	h := handler.NewUserHandler(mockUserService)
	e := echo.New()

	t.Run("should success", func(t *testing.T) {
		reqBody := handler.UserRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockUser := models.User{
			Id:    1,
			Email: reqBody.Email,
			Phone: reqBody.Phone,
		}

		mockUserService.EXPECT().
			LoginUser(reqBody.Email, reqBody.Phone, reqBody.Password).
			Return(&mockUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.LoginUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should bad request when request invalid", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email": 1,
		}
		reqJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.LoginUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should unauthorized when user not found", func(t *testing.T) {
		reqBody := handler.UserRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockUser := models.User{}

		mockUserService.EXPECT().
			LoginUser(reqBody.Email, reqBody.Phone, reqBody.Password).
			Return(&mockUser, errors.New("user not found"))

		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.LoginUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("should bad request when request login incorrect", func(t *testing.T) {
		reqBody := handler.UserRequest{
			Email:    "test@example.com",
			Phone:    "1234567890",
			Password: "password123",
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockUser := models.User{}

		mockUserService.EXPECT().
			LoginUser(reqBody.Email, reqBody.Phone, reqBody.Password).
			Return(&mockUser, errors.New("failed"))

		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.LoginUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
