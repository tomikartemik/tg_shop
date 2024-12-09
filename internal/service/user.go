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

func (s *UserService) CreateUser(id int, username string) (model.User, error) {
	newUser := model.User{
		TelegramID: id,
		Username:   username,
	}
	user, err := s.repo.CreateUser(newUser)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	return s.repo.GetUserById(id)
}
