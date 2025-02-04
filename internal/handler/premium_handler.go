package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_shop/internal/service"
)

type PremiumHandler struct {
	service *service.PremiumService
	bot     *tgbotapi.BotAPI
}

func NewPremiumHandler(service *service.PremiumService, bot *tgbotapi.BotAPI) *PremiumHandler {
	return &PremiumHandler{
		service: service,
		bot:     bot,
	}
}

func (h *PremiumHandler) NotifyPremiumUsers() {
	expiresInThreeDays, expired, err := h.service.GetPremiumInfo()
	if err != nil {
		log.Printf("Error retrieving premium information: %v", err)
		return
	}

	for _, user := range expiresInThreeDays {
		msg := tgbotapi.NewMessage(int64(user.TelegramID),
			"⚠️ Your premium subscription will expire in 3 days. Contact the administrator to renew it.")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending expiration notification to user %d: %v", user.TelegramID, err)
		}
	}

	for _, user := range expired {
		msg := tgbotapi.NewMessage(int64(user.TelegramID),
			"❌ Your premium subscription has expired. You are now limited to 3 ads. To renew, contact the administrator.")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending expiration notification to user %d: %v", user.TelegramID, err)
		}
	}
}
