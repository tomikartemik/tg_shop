package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (repo *InvoiceRepository) CreateInvoice(TelegramID int, amount float64) (int, error) {
	newInvoice := model.Invoice{
		TelegramID: TelegramID,
		Amount:     amount,
		Status:     "Processing",
	}
	result := repo.db.Create(&newInvoice)
	if result.Error != nil {
		return 0, result.Error
	}
	return newInvoice.ID, nil
}
