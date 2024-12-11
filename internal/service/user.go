package service

import (
	"errors"
	"gorm.io/gorm"
	"log"
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

func (s *UserService) CreateOrUpdateUser(user model.User) (model.User, error) {
	existingUser, err := s.repo.GetUserById(user.TelegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found, creating new user with TelegramID: %d", user.TelegramID)
			createdUser, createErr := s.repo.CreateUser(user)
			if createErr != nil {
				log.Printf("Error creating user: %v", createErr)
				return model.User{}, createErr
			}
			return createdUser, nil
		}

		log.Printf("Unexpected error fetching user: %v", err)
		return model.User{}, err
	}

	log.Printf("User found, updating user with TelegramID: %d", existingUser.TelegramID)
	existingUser.Language = user.Language
	if user.Username != "" {
		existingUser.Username = user.Username
	}

	updatedUser, updateErr := s.repo.UpdateUser(existingUser)
	if updateErr != nil {
		log.Printf("Error updating user: %v", updateErr)
		return model.User{}, updateErr
	}

	return updatedUser, nil
}

func (s *UserService) UpdateUserName(telegramID int, username string) error {
	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return err
	}

	user.Username = username
	_, updateErr := s.repo.UpdateUser(user)
	return updateErr
}
