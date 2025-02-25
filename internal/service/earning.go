package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_shop/internal/repository"
)

type EarningService struct {
	repo     repository.Earning
	repoUser repository.User
	bot      *tgbotapi.BotAPI
}

func NewEarningService(repo repository.Earning, repoUser repository.User, bot *tgbotapi.BotAPI) *EarningService {
	return &EarningService{
		repo:     repo,
		repoUser: repoUser,
		bot:      bot,
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

		if err = s.repoUser.IncrementSalesAmount(earning.SellerID); err != nil {
			return fmt.Errorf("failed to increment sales amount: %w", err)
		}

		if err = s.repo.MarkAsProcessed(&earning); err != nil {
			return fmt.Errorf("failed to mark earning as processed: %w", err)
		}

		chatID := earning.BuyerID
		msg := tgbotapi.NewMessage(int64(chatID), fmt.Sprintf("Please, rate seller (ID: %d):", earning.SellerID))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("1", fmt.Sprintf("rate_%d_1", earning.SellerID)),
				tgbotapi.NewInlineKeyboardButtonData("2", fmt.Sprintf("rate_%d_2", earning.SellerID)),
				tgbotapi.NewInlineKeyboardButtonData("3", fmt.Sprintf("rate_%d_3", earning.SellerID)),
				tgbotapi.NewInlineKeyboardButtonData("4", fmt.Sprintf("rate_%d_4", earning.SellerID)),
				tgbotapi.NewInlineKeyboardButtonData("5", fmt.Sprintf("rate_%d_5", earning.SellerID)),
			),
		)

		if _, err := s.bot.Send(msg); err != nil {
			log.Printf("Failed to send rating request: %v", err)
		}
	}

	return nil
}
