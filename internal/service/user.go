package service

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type UserService struct {
	repo         repository.User
	repoAd       repository.Ad
	repoCategory repository.Category
}

func NewUserService(repo repository.User, repoAd repository.Ad, repoCategory repository.Category) *UserService {
	return &UserService{
		repo:         repo,
		repoAd:       repoAd,
		repoCategory: repoCategory,
	}
}

func (s *UserService) CreateUser(id int, user model.User) (model.User, error) {
	newUser, err := s.repo.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	return s.repo.GetUserById(id)
}

func (s *UserService) CreateOrUpdateUser(user model.User) (model.User, error) {
	existingUser, err := s.repo.GetUserById(user.TelegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found, creating new user with TelegramID: %d", user.TelegramID)
			createdUser, createErr := s.repo.CreateUser(user)
			if createErr != nil {
				log.Printf("Error creating user: %v", createErr)
				return model.User{}, createErr
			}
			return createdUser, nil
		}

		log.Printf("Unexpected error fetching user: %v", err)
		return model.User{}, err
	}

	log.Printf("User found, updating user with TelegramID: %d", existingUser.TelegramID)
	existingUser.Language = user.Language
	if user.Username != "" {
		existingUser.Username = user.Username
	}

	updatedUser, updateErr := s.repo.UpdateUser(existingUser)
	if updateErr != nil {
		log.Printf("Error updating user: %v", updateErr)
		return model.User{}, updateErr
	}

	return updatedUser, nil
}

func (s *UserService) UpdateUserName(telegramID int, username string) error {
	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return err
	}

	user.Username = username
	_, updateErr := s.repo.UpdateUser(user)
	return updateErr
}

func (s *UserService) GetUserAsSellerByID(telegramIDStr string) (model.UserAsSeller, error) {
	userAsSeller := model.UserAsSeller{}

	telegramID, err := strconv.Atoi(telegramIDStr)

	if err != nil {
		return model.UserAsSeller{}, err
	}

	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return userAsSeller, err
	}

	adsShortInfo, err := s.convertAdsToAdsShortInfo(user.Ads)

	if err != nil {
		return userAsSeller, err
	}

	userAsSeller = model.UserAsSeller{
		TelegramID:   user.TelegramID,
		Username:     user.Username,
		Ads:          adsShortInfo,
		Rating:       user.Rating,
		ReviewNumber: user.ReviewNumber,
	}

	return userAsSeller, nil
}

func (s *UserService) convertAdsToAdsShortInfo(ads []model.Ad) ([]model.AdShortInfo, error) {
	adsShortInfo := []model.AdShortInfo{}

	for _, ad := range ads {
		seller, err := s.repo.GetUserById(ad.SellerID)

		if err != nil {
			return adsShortInfo, err
		}

		category, err := s.repoCategory.GetCategoryById(ad.CategoryID)

		if err != nil {
			return adsShortInfo, err
		}

		adsShortInfo = append(adsShortInfo, model.AdShortInfo{
			ID:                 ad.ID,
			Title:              ad.Title,
			Price:              ad.Price,
			PhotoURL:           ad.PhotoURL,
			CategoryID:         category.ID,
			CategoryName:       category.Name,
			SellerID:           seller.TelegramID,
			SellerName:         seller.Username,
			SellerRating:       seller.Rating,
			SellerReviewNumber: seller.ReviewNumber,
			Stock:              ad.Stock,
		})
	}

	return adsShortInfo, nil
}
