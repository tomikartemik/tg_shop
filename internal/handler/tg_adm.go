package handler

//import (
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"log"
//	"strconv"
//	"strings"
//)
//
//func (h *AdminHandler) HandleAdminInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	messageText := strings.TrimSpace(update.Message.Text)
//	chatID := update.Message.Chat.ID
//
//	switch messageText {
//	case "/start":
//		h.handleAdminStart(bot, chatID)
//
//	case "📢 Broadcast Message":
//		h.userStates[chatID] = "broadcasting_message"
//		msg := tgbotapi.NewMessage(chatID, "Enter the message to broadcast:")
//		bot.Send(msg)
//
//	case "🔍 Search User":
//		h.userStates[chatID] = "searching_user"
//		msg := tgbotapi.NewMessage(chatID, "Enter the user ID to search:")
//		bot.Send(msg)
//
//	default:
//		if state, exists := h.userStates[chatID]; exists {
//			switch state {
//			case "broadcasting_message":
//				h.handleBroadcast(bot, chatID, messageText)
//			case "searching_user":
//				h.handleUserSearch(bot, update, messageText)
//			default:
//				msg := tgbotapi.NewMessage(chatID, "Command not recognized.")
//				bot.Send(msg)
//			}
//		} else {
//			msg := tgbotapi.NewMessage(chatID, "Command not recognized.")
//			bot.Send(msg)
//		}
//	}
//}
//
//func (h *AdminHandler) handleAdminStart(bot *tgbotapi.BotAPI, chatID int64) {
//	isAdmin, err := h.services.User.IsAdmin(int(chatID))
//	if err != nil {
//		log.Printf("Error checking admin status: %v", err)
//		msg := tgbotapi.NewMessage(chatID, "An error occurred. Please try again later.")
//		bot.Send(msg)
//		return
//	}
//
//	if !isAdmin {
//		msg := tgbotapi.NewMessage(chatID, "You are not authorized to use this bot.")
//		bot.Send(msg)
//		return
//	}
//
//	h.sendAdminMainMenu(bot, chatID)
//}
//
//func (h *AdminHandler) sendAdminMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
//	menuMessage := "Welcome to the Admin Panel. Choose an option below:"
//	menuKeyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("🔍 Search User"),
//			tgbotapi.NewKeyboardButton("📢 Broadcast Message"),
//		),
//	)
//
//	msg := tgbotapi.NewMessage(chatID, menuMessage)
//	msg.ReplyMarkup = menuKeyboard
//
//	_, err := bot.Send(msg)
//	if err != nil {
//		log.Printf("Error sending admin main menu: %v", err)
//	}
//}
//
//func (h *AdminHandler) handleUserSearch(bot *tgbotapi.BotAPI, update tgbotapi.Update, userIDText string) {
//	chatID := update.Message.Chat.ID
//
//	userID, err := strconv.Atoi(userIDText)
//	if err != nil {
//		msg := tgbotapi.NewMessage(chatID, "Invalid user ID. Please enter a numeric user ID:")
//		bot.Send(msg)
//		return
//	}
//
//	user, err := h.services.User.GetUserById(userID)
//	if err != nil {
//		msg := tgbotapi.NewMessage(chatID, "User not found. Please enter a valid user ID:")
//		bot.Send(msg)
//		return
//	}
//
//	userInfo := fmt.Sprintf(
//		"User Info:\nID: %d\nName: %s\nBalance: %.2f\nRating: %.2f\nPremium: %t",
//		user.TelegramID, user.Username, user.Balance, user.Rating, user.IsPremium,
//	)
//	msg := tgbotapi.NewMessage(chatID, userInfo)
//
//	actionKeyboard := tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("🚫 Block User", fmt.Sprintf("block_%d", userID)),
//			tgbotapi.NewInlineKeyboardButtonData("💸 Change Balance", fmt.Sprintf("change_balance_%d", userID)),
//		),
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("⭐ Change Rating", fmt.Sprintf("change_rating_%d", userID)),
//			tgbotapi.NewInlineKeyboardButtonData("👑 Grant Premium", fmt.Sprintf("grant_premium_%d", userID)),
//		),
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("✏️ Edit Ads", fmt.Sprintf("edit_ads_%d", userID)),
//		),
//	)
//
//	msg.ReplyMarkup = actionKeyboard
//	bot.Send(msg)
//}
//
//func (h *AdminHandler) handleBroadcast(bot *tgbotapi.BotAPI, chatID int64, messageText string) {
//	err := h.services.User.BroadcastMessage(messageText)
//	if err != nil {
//		log.Printf("Error broadcasting message: %v", err)
//		msg := tgbotapi.NewMessage(chatID, "Failed to broadcast the message. Please try again.")
//		bot.Send(msg)
//		return
//	}
//
//	msg := tgbotapi.NewMessage(chatID, "Message broadcasted successfully!")
//	bot.Send(msg)
//}
