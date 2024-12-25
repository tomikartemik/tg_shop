package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_shop/internal/handler"
)

func InitBot(botToken string) *tgbotapi.BotAPI {
	if botToken == "" {
		log.Fatal("Telegram bot token not provided")
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
		return nil
	}
	return bot
}

func BotProcess(handlers *handler.Handler, bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				handlers.HandleStart(bot, update)
			default:
				handlers.HandleUserInput(bot, update)
			}
		} else if update.CallbackQuery != nil {
			handlers.HandleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}
