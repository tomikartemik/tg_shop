package repository

import (
	"gorm.io/gorm"
	"tg_shop/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetCategoryList() ([]model.Category, error) {
	categories := []model.Category{}
	err := repo.db.Find(&categories).Error
	return categories, err
}

func (repo *CategoryRepository) GetCategoryById(categoryID int) (model.Category, error) {
	category := model.Category{}
	err := repo.db.First(&category, "id=?", categoryID).Error
	return category, err
}
