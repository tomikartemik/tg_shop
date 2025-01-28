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

func (repo *InvoiceRepository) ChangeStatus(id int, status string) error {
	var invoice model.Invoice
	if err := repo.db.First(&invoice, id).Error; err != nil {
		return err
	}
	invoice.Status = status
	return repo.db.Save(&invoice).Error
}

func (repo *InvoiceRepository) GetInvoiceByID(id int) (model.Invoice, error) {
	var invoice model.Invoice
	result := repo.db.First(&invoice, id)
	return invoice, result.Error
}
