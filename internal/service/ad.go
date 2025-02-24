package service

import (
	"fmt"
	"strconv"
	"strings"
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

func (s *AdService) CreateAd(ad model.Ad) (model.Ad, error) {
	createdAd, err := s.repo.CreateAd(ad)
	if err != nil {
		return model.Ad{}, err
	}

	return createdAd, nil
}

func (s *AdService) GetAdList(categoryIDStr string) ([]model.AdShortInfo, error) {
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

func (s *AdService) GetAdBySellerID(sellerIDStr string) (model.AdInfo, error) {
	id, err := strconv.Atoi(sellerIDStr)

	if err != nil {
		return model.AdInfo{}, err
	}

	ad, err := s.repo.GetAdBySellerId(id)

	if err != nil {
		return model.AdInfo{}, err
	}

	adInfo, err := s.convertAdToAdInfo(ad)

	return adInfo, nil
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

	paragraphs := strings.Split(ad.Description, `\n`)

	// Объединяем абзацы с разделителем \n
	formattedDescription := strings.Join(paragraphs, `\\n`)

	ad.Description = formattedDescription

	adInfo, err := s.convertAdToAdInfo(ad)

	return adInfo, nil
}

func (s *AdService) GetAdsByUserID(userID int) ([]model.Ad, error) {
	ads, err := s.repo.GetAdsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (s *AdService) EditAd(adID int, updatedAd model.Ad) error {
	ad, err := s.repo.GetAdById(adID)
	if err != nil {
		return fmt.Errorf("failed to fetch ad: %w", err)
	}

	ad.Title = updatedAd.Title
	ad.Description = updatedAd.Description
	ad.Price = updatedAd.Price
	ad.Stock = updatedAd.Stock
	ad.CategoryID = updatedAd.CategoryID
	ad.PhotoURL = updatedAd.PhotoURL
	ad.Files = updatedAd.Files

	_, updateErr := s.repo.UpdateAd(ad)
	if updateErr != nil {
		return fmt.Errorf("failed to edit ad: %w", updateErr)
	}

	return nil
}

func (s *AdService) DeleteAd(adID int) error {
	return s.repo.UpdateAdStatus(adID, "Deleted")
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

func (s *AdService) ApproveAd(adID int) error {
	return s.repo.UpdateAdStatus(adID, "Enabled")
}

func (s *AdService) RejectAd(adID int) error {
	return s.repo.UpdateAdStatus(adID, "Rejected")
}

func (s *AdService) GetAdByIDTg(adID int) (model.Ad, error) {
	return s.repo.GetAdByIDTg(adID)
}
