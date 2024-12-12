package service

import (
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetCategoryList() ([]model.Category, error) {
	return s.repo.GetCategoryList()
}
