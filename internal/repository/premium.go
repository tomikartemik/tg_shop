package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
	"time"
)

type PremiumRepository struct {
	db *gorm.DB
}

func NewPremiumRepository(db *gorm.DB) *PremiumRepository {
	return &PremiumRepository{db: db}
}

func (repo *PremiumRepository) GetExpiredPremiums() ([]model.User, []model.User, error) {
	var expiresInThreeDays []model.User
	var expired []model.User
	err := repo.db.Where("expire_premium <= ? AND is_premium IS TRUE", time.Now().Add(72*time.Hour)).
		Order("id ASC").
		Find(&expiresInThreeDays).
		Error

	if err != nil {
		return nil, nil, err
	}

	err = repo.db.Where("expire_premium >= ? AND is_premium IS TRUE", time.Now()).
		Order("id ASC").
		Find(&expiresInThreeDays).
		Error

	if err != nil {
		return nil, nil, err
	}

	return expiresInThreeDays, expired, err
}

// ВОТ ЭТУ ФУНКЦИЮ ВОЗМОЖНО НАДО БУДЕТ ДОПИЛИТЬ
func (repo *PremiumRepository) ResetPremiums(users []model.User) {
	for _, user := range users {
		user.IsPremium = false
		repo.db.Save(&user)
	}
}
