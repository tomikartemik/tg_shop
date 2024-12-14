package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func (h *AdminHandler) HandleAdminStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	adminID := update.Message.From.ID

	isAdmin, err := h.services.User.IsAdmin(int(adminID))
	if err != nil || !isAdmin {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied. You are not an administrator.")
		bot.Send(msg)
		return
	}

	h.sendAdminMainMenu(bot, update.Message.Chat.ID)
}

func (h *AdminHandler) sendAdminMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	menuMessage := "Welcome to the Admin Panel. Choose an option below:"
	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔍 Work with User"),
			tgbotapi.NewKeyboardButton("📢 Broadcast Message"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, menuMessage)
	msg.ReplyMarkup = menuKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending admin main menu: %v", err)
	}
}

func (h *AdminHandler) HandleAdminInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	messageText := update.Message.Text

	switch messageText {
	case "🔍 Work with User":
		h.userStates[chatID] = "searching_user"
		msg := tgbotapi.NewMessage(chatID, "Enter the user ID to search:")
		bot.Send(msg)

	case "📢 Broadcast Message":
		h.userStates[chatID] = "broadcasting_message"
		msg := tgbotapi.NewMessage(chatID, "Enter the message to broadcast:")
		bot.Send(msg)

	default:
		h.sendUnknownCommand(bot, chatID)
	}
}

func (h *AdminHandler) handleUserSearch(bot *tgbotapi.BotAPI, chatID int64, userIDText string) {
	userID, err := strconv.Atoi(userIDText)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid user ID. Please enter a numeric user ID:")
		bot.Send(msg)
		return
	}

	user, err := h.services.User.GetUserById(userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "User not found. Please enter a valid user ID:")
		bot.Send(msg)
		return
	}

	ads, err := h.services.Ad.GetAdsByUserID(userID)
	if err != nil || len(ads) == 0 {
		msg := tgbotapi.NewMessage(chatID, "This user has no ads.")
		bot.Send(msg)
		return
	}

	userInfo := fmt.Sprintf(
		"User Info:\nID: %d\nName: %s\nBalance: %.2f\nRating: %.2f\nPremium: %t\n\nAds:",
		user.TelegramID, user.Username, user.Balance, user.Rating, user.IsPremium,
	)

	for _, ad := range ads {
		userInfo += fmt.Sprintf("\nID: %d, Title: %s, Price: %.2f", ad.ID, ad.Title, ad.Price)
	}

	msg := tgbotapi.NewMessage(chatID, userInfo)

	actionKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫 Block User", fmt.Sprintf("block_%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("💸 Change Balance", fmt.Sprintf("change_balance_%d", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐ Change Rating", fmt.Sprintf("change_rating_%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("👑 Grant Premium", fmt.Sprintf("grant_premium_%d", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Delete Ad", fmt.Sprintf("delete_ad_%d", userID)),
		),
	)

	msg.ReplyMarkup = actionKeyboard
	bot.Send(msg)
}

func (h *AdminHandler) HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID

	switch {
	case strings.HasPrefix(data, "block_"):
		userID, _ := strconv.Atoi(strings.TrimPrefix(data, "block_"))
		h.handleBlockUser(bot, chatID, userID)

	case strings.HasPrefix(data, "change_balance_"):
		userID := strings.TrimPrefix(data, "change_balance_")
		h.userStates[chatID] = fmt.Sprintf("changing_balance_%s", userID)
		msg := tgbotapi.NewMessage(chatID, "Enter the new balance:")
		bot.Send(msg)

	case strings.HasPrefix(data, "change_rating_"):
		userID := strings.TrimPrefix(data, "change_rating_")
		h.userStates[chatID] = fmt.Sprintf("changing_rating_%s", userID)
		msg := tgbotapi.NewMessage(chatID, "Enter the new rating (0-5):")
		bot.Send(msg)

	case strings.HasPrefix(data, "grant_premium_"):
		userID, _ := strconv.Atoi(strings.TrimPrefix(data, "grant_premium_"))
		h.handleGrantPremium(bot, chatID, userID)

	case strings.HasPrefix(data, "delete_ad_"):
		adID := strings.TrimPrefix(data, "delete_ad_")
		h.userStates[chatID] = fmt.Sprintf("deleting_ad_%s", adID)
		msg := tgbotapi.NewMessage(chatID, "Enter the ad ID to delete:")
		bot.Send(msg)

	default:
		msg := tgbotapi.NewMessage(chatID, "Unknown command.")
		bot.Send(msg)
	}

	callbackResponse := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Request(callbackResponse)
}

func (h *AdminHandler) handleGrantPremium(bot *tgbotapi.BotAPI, chatID int64, userID int) {
	err := h.services.User.GrantPremium(userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to grant premium status. Please try again.")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Premium status granted successfully.")
	bot.Send(msg)
}

func (h *AdminHandler) handleBlockUser(bot *tgbotapi.BotAPI, chatID int64, userID int) {
	err := h.services.User.BlockUser(userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to block user. Please try again.")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "User has been blocked successfully.")
	bot.Send(msg)
}
func (h *AdminHandler) sendUnknownCommand(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Command not recognized.")
	bot.Send(msg)
}
