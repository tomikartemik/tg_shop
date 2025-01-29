package service

import (
	"fmt"
	"tg_shop/internal/repository"
)

type PremiumService struct {
	repo repository.Premium
}

func NewPremiumService(repo repository.Premium) *PremiumService {
	return &PremiumService{
		repo: repo,
	}
}

// ВОТ ТУТ ПОЛУЧАЕМ ИНФУ
// expiresInThreeDays - те, у кого премиум истекает меньше, чем через 3 дня, им кидаем сообщение
// expiresInThreeDays - те, у кого премиум истек, сам знаешь, что с ними делать
// return просто прописал, пока не очень понимаю, что вообще возвращать надо будет
func (s *PremiumService) GetPremiumInfo() error {
	expiresInThreeDays, expired, err := s.repo.GetExpiredPremiums()

	fmt.Println(expiresInThreeDays)
	fmt.Println(expired)
	//строчки выше, просто чтоб запускалось, надо как-то переменные заюзать

	if err != nil {
		return err
	}
	return err
}
