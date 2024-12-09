package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
)

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID

	user, err := h.services.GetUserById(int(telegramID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Пожалуйста, введите своё имя:")
			h.pendingUsernames[telegramID] = true // Добавляем в ожидание
			_, sendErr := bot.Send(msg)
			if sendErr != nil {
				log.Printf("Error sending message: %v", sendErr)
			}
			return
		}
		log.Printf("Error checking user: %v", err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s! Ваш баланс: %.2f", user.Username, user.Balance))
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		log.Printf("Error sending message: %v", sendErr)
	}
}

func (h *Handler) HandleUserInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	messageText := update.Message.Text

	// Проверяем, ждём ли мы ввода имени от пользователя
	if h.pendingUsernames[telegramID] {
		// Сохраняем пользователя в базе
		newUser, err := h.services.CreateUser(int(telegramID), messageText)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при создании пользователя. Попробуйте снова.")
			bot.Send(msg)
			return
		}

		// Уведомляем об успешном создании
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваш профиль создан!\nВаше имя: %s\nБаланс: %.2f", newUser.Username, newUser.Balance))
		bot.Send(msg)

		// Удаляем пользователя из списка ожидающих
		delete(h.pendingUsernames, telegramID)
	} else {
		// Если мы не ожидаем имя, просто игнорируем сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не распознана.")
		bot.Send(msg)
	}
}
