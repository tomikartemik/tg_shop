package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("error with loading data from env: %s", err.Error())
	}

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	for update := range updates {
		if update.Message != nil {
			chatID := tu.ID(update.Message.Chat.ID)

			_, _ = bot.CopyMessage(
				tu.CopyMessage(
					chatID,
					chatID,
					update.Message.MessageID,
				),
			)
		}
	}
}
