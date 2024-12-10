package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"tg_shop/internal/model"
)

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID

	h.userStates[telegramID] = "language"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите язык:")
	languageKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇷🇺Русский"),
			tgbotapi.NewKeyboardButton("🇺🇸English"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🇪🇸Spanish"),
			tgbotapi.NewKeyboardButton("🇩🇪Deutsch"),
		),
	)
	msg.ReplyMarkup = languageKeyboard
	bot.Send(msg)
}

func (h *Handler) HandleUserInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	messageText := strings.TrimSpace(update.Message.Text)

	log.Printf("User %d state: %s", telegramID, h.userStates[telegramID])
	log.Printf("Received message: %s", messageText)
	log.Printf(h.userStates[telegramID])

	if h.userStates[telegramID] == "language" {
		var language string
		switch messageText {
		case "🇷🇺Русский":
			language = "ru"
		case "🇺🇸English":
			language = "en"
		case "🇪🇸Spanish":
			language = "es"
		case "🇩🇪Deutsch":
			language = "de"
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите язык из предложенных вариантов.")
			bot.Send(msg)
			return
		}

		delete(h.userStates, telegramID)
		h.userStates[telegramID] = "username"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, h.getLocalizedMessage(language, "Теперь введите своё имя:"))
		bot.Send(msg)
		return
	}

	if h.userStates[telegramID] == "username" {
		delete(h.userStates, telegramID)
		return
	}

	if language != "" {
		newUser := model.User{
			TelegramID: int(telegramID),
			Username:   messageText,
			Language:   language,
		}

		savedUser, err := h.services.CreateUser(newUser.TelegramID, newUser)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при создании пользователя. Попробуйте снова.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Ваш профиль создан!\nВаше имя: %s\nЯзык: %s\nБаланс: %.2f",
			savedUser.Username, savedUser.Language, savedUser.Balance,
		))
		bot.Send(msg)

		delete(h.userStates, telegramID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не распознана.")
	bot.Send(msg)
}

func (h *Handler) getLocalizedMessage(language, defaultMessage string) string {
	messages := map[string]map[string]string{
		"ru": {"Теперь введите своё имя:": "Теперь введите своё имя:"},
		"en": {"Теперь введите своё имя:": "Please enter your name:"},
		"es": {"Теперь введите своё имя:": "Por favor, introduzca su nombre:"},
		"de": {"Теперь введите своё имя:": "Bitte geben Sie Ihren Namen ein:"},
	}

	if localized, ok := messages[language][defaultMessage]; ok {
		return localized
	}
	return defaultMessage
}
