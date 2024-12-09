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
	err := repo.db.Create(user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserById(id int) (model.User, error) {
	var user model.User
	err := repo.db.Where("telegram_id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
