package repository

import (
	"fmt"
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type AdRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) *AdRepository {
	return &AdRepository{db: db}
}

func (repo *AdRepository) CreateAd(ad model.Ad) (model.Ad, error) {
	err := repo.db.Create(&ad).Error
	if err != nil {
		return model.Ad{}, err
	}
	return ad, nil
}

func (repo *AdRepository) GetAllAds() ([]model.Ad, error) {
	var ads []model.Ad
	err := repo.db.Find(&ads).Error
	if err != nil {
		return ads, err
	}
	return ads, nil
}

func (repo *AdRepository) GetAdListByCategory(categoryID int) ([]model.Ad, error) {
	var ads []model.Ad
	err := repo.db.Where("category_id = ?", categoryID).Find(&ads).Error
	//err := repo.db.Model(model.Ad{}).Find(&ads).Error
	if err != nil {
		return []model.Ad{}, err
	}
	return ads, nil
}

func (repo *AdRepository) GetAdBySellerId(id int) (model.Ad, error) {
	ad := model.Ad{}
	err := repo.db.Where("seller_id = ?", id).First(&ad).Error
	if err != nil {
		return model.Ad{}, err
	}
	return ad, nil
}

func (repo *AdRepository) GetAdById(id int) (model.Ad, error) {
	ad := model.Ad{}
	err := repo.db.Where("id = ?", id).First(&ad).Error
	if err != nil {
		return model.Ad{}, err
	}
	return ad, nil
}

func (repo *AdRepository) GetAdsByUserID(userID int) ([]model.Ad, error) {
	var ads []model.Ad
	err := repo.db.Where("seller_id = ?", userID).Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (repo *AdRepository) UpdateAd(ad model.Ad) (model.Ad, error) {
	err := repo.db.Model(&ad).Updates(ad).Error
	if err != nil {
		return model.Ad{}, err
	}
	return ad, nil
}

func (repo *AdRepository) DeleteAd(adID int) error {
	result := repo.db.Delete(&model.Ad{}, adID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no ad found with ID: %d", adID)
	}

	return nil
}
