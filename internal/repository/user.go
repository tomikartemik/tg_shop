package repository

import (
	"fmt"
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user model.User) (model.User, error) {
	err := repo.db.Create(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserById(id int) (model.User, error) {
	var user model.User
	err := repo.db.Where("telegram_id = ?", id).Preload("Ads").Preload("Purchased").First(&user).Error
	return user, err
}

func (repo *UserRepository) UpdateUser(user model.User) (model.User, error) {
	var existingUser model.User
	err := repo.db.Where("telegram_id = ?", user.TelegramID).First(&existingUser).Error
	if err != nil {
		return model.User{}, err
	}

	err = repo.db.Model(&existingUser).Updates(user).Error
	if err != nil {
		return model.User{}, err
	}

	return existingUser, nil
}

func (repo *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := repo.db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (repo *UserRepository) SearchUsers(query string) ([]model.User, error) {
	var users []model.User
	err := repo.db.Where("username ILIKE ?", "%"+query+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
