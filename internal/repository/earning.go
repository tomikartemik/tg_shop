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
	return r.db.Save(earning).Error
}
