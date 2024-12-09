package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"tg_shop/internal/handler"
)

func BotProcess(handlers *handler.Handler) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Telegram bot token not provided")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				handlers.HandleStart(bot, update)
			default:
				log.Printf("Unhandled message: %s", update.Message.Text)
			}
		}
	}
}
