package service

import (
	"fmt"
	"strconv"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type AdService struct {
	repo         repository.Ad
	repoUser     repository.User
	repoCategory repository.Category
}

func NewAdService(repo repository.Ad, repoUser repository.User, repoCategory repository.Category) *AdService {
	return &AdService{
		repo:         repo,
		repoUser:     repoUser,
		repoCategory: repoCategory,
	}
}

func (s *AdService) CreateAd() {
	// Сам пропишешь, не знаю, как тебе удобнее будет
}

func (s *AdService) GetAdList(categoryIDStr string) ([]model.AdShortInfo, error) {
	fmt.Println(categoryIDStr)
	adsShortInfo := []model.AdShortInfo{}

	if categoryIDStr == "" {
		ads, err := s.repo.GetAllAds()

		if err != nil {
			return adsShortInfo, err
		}

		return s.convertAdsToAdsShortInfo(ads)
	}

	categoryID, err := strconv.Atoi(categoryIDStr)

	if err != nil {
		return adsShortInfo, err
	}

	ads, err := s.repo.GetAdListByCategory(categoryID)
	if err != nil {
		return adsShortInfo, err
	}

	return s.convertAdsToAdsShortInfo(ads)
}

func (s *AdService) GetAdByID(idStr string) (model.AdInfo, error) {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return model.AdInfo{}, err
	}

	ad, err := s.repo.GetAdById(id)

	if err != nil {
		return model.AdInfo{}, err
	}

	adInfo, err := s.convertAdToAdInfo(ad)

	return adInfo, nil
}

func (s *AdService) convertAdToAdInfo(ad model.Ad) (model.AdInfo, error) {
	adInfo := model.AdInfo{}

	seller, err := s.repoUser.GetUserById(ad.SellerID)

	if err != nil {
		return adInfo, err
	}

	category, err := s.repoCategory.GetCategoryById(ad.CategoryID)

	if err != nil {
		return adInfo, err
	}

	adInfo = model.AdInfo{
		ID:                 ad.ID,
		Title:              ad.Title,
		Description:        ad.Description,
		Price:              ad.Price,
		Files:              ad.Files,
		PhotoURL:           ad.PhotoURL,
		CategoryID:         category.ID,
		CategoryName:       category.Name,
		SellerID:           seller.TelegramID,
		SellerName:         seller.Username,
		SellerRating:       seller.Rating,
		SellerReviewNumber: seller.ReviewNumber,
		Stock:              ad.Stock,
	}

	return adInfo, nil
}

func (s *AdService) convertAdsToAdsShortInfo(ads []model.Ad) ([]model.AdShortInfo, error) {
	adsShortInfo := []model.AdShortInfo{}

	for _, ad := range ads {
		seller, err := s.repoUser.GetUserById(ad.SellerID)

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
