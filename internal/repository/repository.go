package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}

type User interface {
	CreateUser(user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
}
