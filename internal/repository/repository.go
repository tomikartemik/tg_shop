package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type Repository struct {
	User
	Ad
	Category
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Ad:       NewAdRepository(db),
		Category: NewCategoryRepository(db),
	}
}

type User interface {
	CreateUser(user model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
}

type Ad interface {
	CreateAd(ad model.Ad) (model.Ad, error)
	GetAdListByCategory(categoryID int) ([]model.Ad, error)
	GetAllAds() ([]model.Ad, error)
	GetAdBySellerId(id int) (model.Ad, error)
	GetAdsByUserID(userID int) ([]model.Ad, error)
	GetAdById(id int) (model.Ad, error)
}

type Category interface {
	GetCategoryList() ([]model.Category, error)
	GetCategoryById(categoryID int) (model.Category, error)
}
