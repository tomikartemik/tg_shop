package repository

import (
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
		return user, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserById(id int) (model.User, error) {
	var user model.User
	err := repo.db.Where("telegram_id = ?", id).First(&user).Error
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
