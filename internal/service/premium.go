package service

import (
	"log"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type PremiumService struct {
	repo   repository.Premium
	repoAd repository.Ad
}

func NewPremiumService(repo repository.Premium, repoAd repository.Ad) *PremiumService {
	return &PremiumService{
		repo:   repo,
		repoAd: repoAd,
	}
}

func (s *PremiumService) GetPremiumInfo() ([]model.User, []model.User, error) {
	expiresInThreeDays, expired, err := s.repo.GetExpiredPremiums()
	if err != nil {
		log.Printf("Ошибка получения информации о премиумах: %v", err)
		return nil, nil, err
	}

	// Если есть истёкшие премиумы, отключаем их
	if len(expired) > 0 {
		err := s.repo.ResetPremiums(expired) // Массовое обновление
		if err != nil {
			log.Printf("Ошибка сброса премиума: %v", err)
			return nil, nil, err
		}
		log.Printf("✅ Отключен премиум у %d пользователей", len(expired))

		for _, user := range expired {
			err := s.repoAd.DisableExcessAds(user.TelegramID)
			if err != nil {
				log.Printf("Ошибка при отключении лишних объявлений у пользователя %d: %v", user.TelegramID, err)
			}
		}
	}
	return expiresInThreeDays, expired, nil
}
