package service

import (
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(id int, user model.User) (model.User, error) {
	newUser, err := s.repo.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	return s.repo.GetUserById(id)
}
