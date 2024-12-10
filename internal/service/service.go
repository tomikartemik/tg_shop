package service

import (
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}

type User interface {
	CreateUser(id int, user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
}
