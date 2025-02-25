package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type Service struct {
	User
	Ad
	Category
	CryptoCloud
	Payout
	Earning
	Premium
}

func NewService(repos *repository.Repository, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		User:        NewUserService(repos.User, repos.Ad, repos.Category, repos.Earning, bot),
		Ad:          NewAdService(repos.Ad, repos.User, repos.Category),
		Category:    NewCategoryService(repos.Category),
		CryptoCloud: NewCryptoCloudService(repos.User, repos.Invoice, bot),
		Payout:      NewPayoutService(repos.Payout),
		Earning:     NewEarningService(repos.Earning, repos.User, bot),
		Premium:     NewPremiumService(repos.Premium, repos.Ad),
	}
}

type User interface {
	CreateUser(id int, user model.User) (model.User, error)
	GetUserInfoById(id int) (model.UserInfo, error)
	CreateOrUpdateUser(user model.User) (model.User, error)
	GetUserAsSellerByID(telegramIDStr string) (model.UserAsSeller, error)
	IsAdmin(userID int) (bool, error)
	BroadcastMessage(message string) error
	BlockUser(userID int) error
	GrantPremium(userID int) error
	ChangeBalance(userID int, newBalance float64) error
	ChangeRating(sellerID int, review int) error
	GetUserByUsername(username string) (model.User, error)
	SearchUsers(query string) ([]model.UserInfo, error)
	Purchase(request model.PurchaseRequest) error
	GetUserById(id int) (model.User, error)
	ChangeRatingAdm(userID int, newRating float64) error
	BroadcastAboutDelete(sellerID int, message string) error
}

type Ad interface {
	GetAdList(categoryIDStr string) ([]model.AdShortInfo, error)
	GetAdBySellerID(idStr string) (model.AdInfo, error)
	CreateAd(ad model.Ad) (model.Ad, error)
	GetAdsByUserID(userID int) ([]model.Ad, error)
	GetAdByID(idStr string) (model.AdInfo, error)
	EditAd(adID int, updatedAd model.Ad) error
	DeleteAd(adID int) error
	ApproveAd(adID int) error
	RejectAd(adID int) error
	GetAdByIDTg(adID int) (model.Ad, error)
}

type Category interface {
	GetCategoryList() ([]model.Category, error)
	GetCategoryById(id int) (model.Category, error)
}

type CryptoCloud interface {
	CreateInvoice(amount float64, telegramID int) (string, error)
	ChangeStatus(id string, status string) error
}

type Payout interface {
	CreatePayoutRequest(telegramID int, amount float64) (int, error)
	ApprovePayoutRequest(requestID int) error
	RejectPayoutRequest(requestID int) error
	GetPayoutByID(telegramID int) (model.PayoutRequest, error)
}

type Earning interface {
	ProcessEarnings() error
}

type Premium interface {
	GetPremiumInfo() ([]model.User, []model.User, error)
}
