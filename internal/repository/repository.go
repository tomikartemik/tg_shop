package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type Repository struct {
	User
	Ad
	Category
	Invoice
	Payout
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Ad:       NewAdRepository(db),
		Category: NewCategoryRepository(db),
		Invoice:  NewInvoiceRepository(db),
		Payout:   NewPayoutRequestRepository(db),
	}
}

type User interface {
	CreateUser(user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
	SearchUsers(query string) ([]model.User, error)
	AddPurchase(userID, adID int) error
	ChangeBalance(userID int, newBlance float64) error
}

type Ad interface {
	CreateAd(ad model.Ad) (model.Ad, error)
	GetAdListByCategory(categoryID int) ([]model.Ad, error)
	GetAllAds() ([]model.Ad, error)
	GetAdBySellerId(id int) (model.Ad, error)
	GetAdsByUserID(userID int) ([]model.Ad, error)
	GetAdById(id int) (model.Ad, error)
	UpdateAd(ad model.Ad) (model.Ad, error)
	DeleteAd(adID int) error // Добавляем метод удаления
	ChangeStock(adID, newStock int) error
	UpdateAdStatus(id int, b bool) error
	GetAdByIDTg(id int) (model.Ad, error)
}

type Category interface {
	GetCategoryList() ([]model.Category, error)
	GetCategoryById(categoryID int) (model.Category, error)
}

type Invoice interface {
	CreateInvoice(TelegramID int, amount float64) (int, error)
}

type Payout interface {
	CreatePayoutRequest(telegramID int, amount float64) (int, error)
	UpdatePayoutStatus(requestID int, status string) error
	GetPayoutByID(telegramID int) (model.PayoutRequest, error)
}
