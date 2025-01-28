package service

import (
	"fmt"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
)

type PayoutService struct {
	repo repository.Payout
}

func NewPayoutService(repo repository.Payout) *PayoutService {
	return &PayoutService{repo: repo}
}

func (s *PayoutService) CreatePayoutRequest(telegramID int, amount float64) (int, error) {
	if amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	return s.repo.CreatePayoutRequest(telegramID, amount)
}

func (s *PayoutService) ApprovePayoutRequest(requestID int) error {
	return s.repo.UpdatePayoutStatus(requestID, "Approved")
}

func (s *PayoutService) RejectPayoutRequest(requestID int) error {
	return s.repo.UpdatePayoutStatus(requestID, "Rejected")
}

func (s *PayoutService) GetPayoutByID(payoutID int) (model.PayoutRequest, error) {
	return s.repo.GetPayoutByID(payoutID)
}
