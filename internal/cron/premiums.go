package cron

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_shop/internal/repository"
	"time"
)

type PremiumsCron struct {
	bot      *tgbotapi.BotAPI
	repoUser repository.User
}

func NewPremiumsCron(bot *tgbotapi.BotAPI, repoUser repository.User) *PremiumsCron {
	return &PremiumsCron{
		bot:      bot,
		repoUser: repoUser,
	}
}

func (c *PremiumsCron) CheckPremiumExpiry() {
	users, err := c.repoUser.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
		return
	}

	for _, user := range users {
		daysLeft := int(user.ExpirePremium.Sub(time.Now()).Hours() / 24)

		if daysLeft == 0 || daysLeft == 3 {
			c.sendPremiumExpiryNotification(int64(user.TelegramID), daysLeft)
		}
	}
}

func (c *PremiumsCron) sendPremiumExpiryNotification(userID int64, daysLeft int) {
	var message string
	if daysLeft == 0 {
		message = "Your premium membership expires today!"
	} else if daysLeft == 3 {
		message = "Your premium membership will expire in 3 days!"
	}

	msg := tgbotapi.NewMessage(userID, message)
	if _, err := c.bot.Send(msg); err != nil {
		log.Printf("Failed to send message to %d: %v", userID, err)
	}
}
