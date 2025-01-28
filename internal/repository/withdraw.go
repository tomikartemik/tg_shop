package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type PayoutRequestRepository struct {
	db *gorm.DB
}

func NewPayoutRequestRepository(db *gorm.DB) *PayoutRequestRepository {
	return &PayoutRequestRepository{db: db}
}

func (repo *PayoutRequestRepository) GetPayoutByID(payoutID int) (model.PayoutRequest, error) {
	var payout model.PayoutRequest
	result := repo.db.Where("id = ?", payoutID).First(&payout)
	return payout, result.Error
}

func (repo *PayoutRequestRepository) CreatePayoutRequest(telegramID int, amount float64) (int, error) {
	newPayoutRequest := model.PayoutRequest{
		TelegramID: telegramID,
		Amount:     amount,
		Status:     "Pending",
	}
	result := repo.db.Create(&newPayoutRequest)
	if result.Error != nil {
		return 0, result.Error
	}
	return newPayoutRequest.ID, nil
}

func (repo *PayoutRequestRepository) UpdatePayoutStatus(requestID int, status string) error {
	result := repo.db.Model(&model.PayoutRequest{}).Where("id = ?", requestID).Update("status", status)
	return result.Error
}
