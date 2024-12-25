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
}

func NewService(repos *repository.Repository, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		User:     NewUserService(repos.User, repos.Ad, repos.Category, bot),
		Ad:       NewAdService(repos.Ad, repos.User, repos.Category),
		Category: NewCategoryService(repos.Category),
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
	ChangeRating(userID int, newRating float64) error
	GetUserByUsername(username string) (model.User, error)
	SearchUsers(query string) ([]model.UserInfo, error)
	Purchase(request model.PurchaseRequest) error
	GetUserById(id int) (model.User, error)
}

type Ad interface {
	GetAdList(categoryIDStr string) ([]model.AdShortInfo, error)
	GetAdBySellerID(idStr string) (model.AdInfo, error)
	CreateAd(ad model.Ad) (model.Ad, error)
	GetAdsByUserID(userID int) ([]model.AdShortInfo, error)
	GetAdByID(idStr string) (model.AdInfo, error)
	EditAd(adID int, updatedAd model.Ad) error
	DeleteAd(adID int) error
}

type Category interface {
	GetCategoryList() ([]model.Category, error)
}
