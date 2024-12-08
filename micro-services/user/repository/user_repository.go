package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"monorepo-ecommerce/micro-services/user/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(user models.User) (*models.User, error)
	GetUserByEmailOrPhone(email string, phone string, password string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.User) (res *models.User, err error) {
	query := `INSERT INTO users (email, phone, password) VALUES (?, ?, ?)`
	data, err := r.db.Exec(query, user.Email, user.Phone, user.Password)
	if err != nil {
		return res, err
	}

	id, _ := data.LastInsertId()
	res = &models.User{
		Id:    id,
		Email: user.Email,
		Phone: user.Phone,
	}

	return res, nil
}

func (r *userRepository) GetUserByEmailOrPhone(email string, phone string, password string) (user *models.User, err error) {
	query := `SELECT id, email, phone, password FROM users WHERE email = ? OR phone = ?`
	row := r.db.QueryRow(query, email, phone)

	var data models.User
	if err := row.Scan(&data.Id, &data.Email, &data.Phone, &data.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	user = &data
	user.Password = ""

	return user, nil
}
