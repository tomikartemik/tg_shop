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

	existingUser, err := h.services.GetUserById(int(telegramID))
	if err == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Добро пожаловать обратно, %s! Ваш текущий язык: %s",
			existingUser.Username, existingUser.Language,
		))
		bot.Send(msg)
		sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Printf("Error checking user existence: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте позже.")
		bot.Send(msg)
		return
	}

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

		log.Printf("Langugage: %s", language)
		newUser := model.User{
			TelegramID: int(telegramID),
			Language:   language,
		}
		log.Printf("User: %s", newUser)
		_, err := h.services.CreateOrUpdateUser(newUser)
		if err != nil {
			log.Printf("Error creating/updating user: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при сохранении языка. Попробуйте снова.")
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
		username := messageText

		user := model.User{
			TelegramID: int(telegramID),
			Username:   username,
		}

		_, err := h.services.CreateOrUpdateUser(user)
		if err != nil {
			log.Printf("Error updating username: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при сохранении имени. Попробуйте снова.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваше имя сохранено: %s", username))
		bot.Send(msg)

		delete(h.userStates, telegramID)
		sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	if h.userStates[telegramID] == "" {
		switch messageText {
		case "📋 Создать объявление":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите данные для нового объявления.")
			bot.Send(msg)
		case "👤 Профиль":
			user, err := h.services.GetUserById(int(telegramID))
			if err != nil {
				log.Printf("Error fetching user profile: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при загрузке профиля.")
				bot.Send(msg)
				return
			}
			profileMessage := fmt.Sprintf("Ваш профиль:\nИмя: %s\nБаланс: %.2f\nРейтинг: %.2f", user.Username, user.Balance, user.Rating)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, profileMessage)
			bot.Send(msg)
		case "📌 Важное":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Важная информация...")
			bot.Send(msg)
		case "💬 Поддержка":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Свяжитесь с нашей поддержкой по вопросам...")
			bot.Send(msg)
		case "🌐 Наши сервисы":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Наши сервисы:\n1. Сервис A\n2. Сервис B...")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите действие из меню.")
			bot.Send(msg)
		}
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

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	menuMessage := "Выберите действие из меню ниже:"
	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Создать объявление"),
			tgbotapi.NewKeyboardButton("👤 Профиль"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📌 Важное"),
			tgbotapi.NewKeyboardButton("💬 Поддержка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🌐 Наши сервисы"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, menuMessage)
	msg.ReplyMarkup = menuKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending main menu: %v", err)
	}
}
