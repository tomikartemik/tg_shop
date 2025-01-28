package service

import (
	"fmt"
	"tg_shop/internal/repository"
)

type EarningService struct {
	repo     repository.Earning
	repoUser repository.User
}

func NewEarningService(repo repository.Earning, repoUser repository.User) *EarningService {
	return &EarningService{
		repo:     repo,
		repoUser: repoUser,
	}
}

func (s *EarningService) ProcessEarnings() error {
	earnings, err := s.repo.GetUnprocessedEarnings()

	if err != nil {
		return fmt.Errorf("failed to get unprocessed earnings: %w", err)
	}

	for _, earning := range earnings {
		seller, err := s.repoUser.GetUserById(earning.SellerID)

		if err != nil {
			return fmt.Errorf("failed to get seller by id: %w", err)
		}

		newBalance := seller.Balance + earning.Amount

		if err = s.repoUser.ChangeBalance(earning.SellerID, newBalance); err != nil {
			return fmt.Errorf("failed to change balance: %w", err)
		}

		newHoldBalance := seller.HoldBalance - earning.Amount

		if err = s.repoUser.ChangeHoldBalance(earning.SellerID, newHoldBalance); err != nil {
			return fmt.Errorf("failed to change balance: %w", err)
		}

		if err = s.repo.MarkAsProcessed(&earning); err != nil {
			return fmt.Errorf("failed to mark earning as processed: %w", err)
		}
	}

	return nil
}
