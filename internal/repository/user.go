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

	err := repo.db.Where("username ILIKE ?", query+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) AddPurchase(userID, adID int) error {
	var user model.User
	var ad model.Ad

	if err := repo.db.First(&user, userID).Error; err != nil {
		return err
	}

	if err := repo.db.Model(&model.Ad{}).First(&ad, adID).Error; err != nil {
		return err
	}

	if err := repo.db.Model(&user).Association("Purchased").Append(&ad); err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) ChangeBalance(userID int, newBalance float64) error {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var user model.User

	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.Balance = newBalance

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserRepository) ChangeHoldBalance(userID int, newBalance float64) error {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var user model.User

	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.HoldBalance = newBalance

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserRepository) IncrementSalesAmount(userID int) error {
	if err := repo.db.Table("users").
		Where("telegram_id = ?", userID).
		Update("sales_amount", gorm.Expr("sales_amount + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}
