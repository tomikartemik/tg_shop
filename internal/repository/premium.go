package repository

import (
	"gorm.io/gorm"
	"log"
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
	now := time.Now()

	// Получаем пользователей, у которых премиум истекает через 3 дня
	err := repo.db.Where("expire_premium BETWEEN ? AND ? AND is_premium = TRUE", now.Add(71*time.Hour), now.Add(72*time.Hour)).
		Order("expire_premium ASC").
		Find(&expiresInThreeDays).
		Error
	if err != nil {
		return nil, nil, err
	}

	err = repo.db.Where("expire_premium < ? AND is_premium = TRUE", now).
		Order("expire_premium ASC").
		Find(&expired).
		Error
	if err != nil {
		return nil, nil, err
	}

	return expiresInThreeDays, expired, nil
}

func (repo *PremiumRepository) ResetPremiums(users []model.User) error {
	for _, user := range users {
		err := repo.db.Model(&model.User{}).
			Where("telegram_id = ?", user.TelegramID).
			Update("is_premium", false).
			Error
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
