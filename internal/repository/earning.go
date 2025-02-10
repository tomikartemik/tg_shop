package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
	"time"
)

type EarningRepository struct {
	db *gorm.DB
}

func NewEarningRepository(db *gorm.DB) *EarningRepository {
	return &EarningRepository{db: db}
}

func (r *EarningRepository) CreateEarning(newEarning model.Earning) error {
	if err := r.db.Create(&newEarning).Error; err != nil {
		return err
	}
	return nil
}

func (r *EarningRepository) GetUnprocessedEarnings() ([]model.Earning, error) {
	var earnings []model.Earning
	err := r.db.Where("created_at <= ? AND processed_at IS NULL", time.Now().Add(-72*time.Hour)).
		Order("id ASC").
		Find(&earnings).
		Error
	return earnings, err
}

func (r *EarningRepository) MarkAsProcessed(earning *model.Earning) error {
	now := time.Now()
	earning.ProcessedAt = &now
	earning.Status = "Processed"
	return r.db.Save(earning).Error
}

func (r *EarningRepository) CountEarningsById(telegramID int) (int, error) {
	var count int64
	if err := r.db.Model(&model.Earning{}).Where("seller_id = ?", telegramID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
