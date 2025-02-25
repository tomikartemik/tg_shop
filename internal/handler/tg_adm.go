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
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📢 Broadcast Message"),
			tgbotapi.NewKeyboardButton("🔎 Find User by Username"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Cancel"),
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

	if messageText == "❌ Cancel" {
		delete(h.userStates, chatID)
		h.sendAdminMainMenu(bot, chatID)
		return
	}

	if state, exists := h.userStates[chatID]; exists {
		switch {
		case state == "searching_user_nickname":
			user, err := h.services.User.GetUserByUsername(messageText)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("User with username '%s' not found.", messageText))
				bot.Send(msg)
				return
			}
			h.handleUserSearch(bot, chatID, strconv.Itoa(user.TelegramID))
		case state == "broadcasting_message":
			h.handleBroadcast(bot, chatID, messageText)
		case state == "searching_user_by_username":
			h.handleFindUserByUsername(bot, chatID, messageText)
			delete(h.userStates, chatID)
		case strings.HasPrefix(state, "changing_balance_"):
			userID, _ := strconv.Atoi(strings.TrimPrefix(state, "changing_balance_"))
			h.handleChangeBalance(bot, chatID, userID, messageText)
		case strings.HasPrefix(state, "changing_rating_"):
			userID, _ := strconv.Atoi(strings.TrimPrefix(state, "changing_rating_"))
			h.handleChangeRating(bot, chatID, userID, messageText)
		case strings.HasPrefix(state, "deleting_ad_"):
			adID, _ := strconv.Atoi(strings.TrimPrefix(state, "deleting_ad_"))
			h.handleDeleteAd(bot, chatID, adID, messageText)
		case strings.HasPrefix(state, "waiting_for_ad_id_"):
			adID, _ := strconv.Atoi(strings.TrimPrefix(state, "deleting_ad_"))

			h.handleDeleteAd(bot, chatID, adID, messageText)

			delete(h.userStates, chatID)
			return
		default:
			h.sendUnknownCommand(bot, chatID)
		}
		delete(h.userStates, chatID)
		return
	}

	switch messageText {
	case "🔍 Work with User":
		h.userStates[chatID] = "searching_user_nickname"
		msg := tgbotapi.NewMessage(chatID, "Enter the user nickname to search:")
		bot.Send(msg)

	case "📢 Broadcast Message":
		h.userStates[chatID] = "broadcasting_message"
		msg := tgbotapi.NewMessage(chatID, "Enter the message to broadcast:")
		bot.Send(msg)

	case "🔎 Find User by Username":
		h.userStates[chatID] = "searching_user_by_username"
		msg := tgbotapi.NewMessage(chatID, "Enter the username of the user:")
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
	userInfo := fmt.Sprintf(
		"User Info:\nID: %d\nName: %s\nBalance: %.2f\nRating: %.2f\nPremium: %t",
		user.TelegramID, user.Username, user.Balance, user.Rating, user.IsPremium,
	)

	if len(ads) > 0 {
		userInfo += "\n\nUser Ads:"
		for _, ad := range ads {
			userInfo += fmt.Sprintf("\n- ID: %d, Title: %s, Price: %.2f, Status: %s", ad.ID, ad.Title, ad.Price, ad.Status)
		}
	} else {
		userInfo += "\n\nNo ads found for this user."
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "cancel"),
		),
	)

	msg.ReplyMarkup = actionKeyboard
	bot.Send(msg)
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

func (h *AdminHandler) handleChangeBalance(bot *tgbotapi.BotAPI, chatID int64, userID int, balanceText string) {
	newBalance, err := strconv.ParseFloat(balanceText, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid balance. Please enter a numeric value.")
		bot.Send(msg)
		return
	}

	err = h.services.User.ChangeBalance(userID, newBalance)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to change balance.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Balance updated successfully.")
	bot.Send(msg)
}

func (h *AdminHandler) handleChangeRating(bot *tgbotapi.BotAPI, chatID int64, userID int, ratingText string) {
	newRating, err := strconv.ParseFloat(ratingText, 64)
	if err != nil || newRating < 0 || newRating > 5 {
		msg := tgbotapi.NewMessage(chatID, "Invalid rating. Please enter a value between 0 and 5.")
		bot.Send(msg)
		return
	}

	err = h.services.User.ChangeRatingAdm(userID, newRating)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to change rating.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Rating updated successfully.")
	bot.Send(msg)
}

func (h *AdminHandler) handleDeleteAd(bot *tgbotapi.BotAPI, chatID int64, adID int, messageText string) {
	ad, err := h.services.Ad.GetAdByID(messageText)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ad not found. Please check the ID and try again.")
		bot.Send(msg)
		return
	}

	if ad.SellerID == 0 {
		msg := tgbotapi.NewMessage(chatID, "Seller ID is invalid for this ad.")
		bot.Send(msg)
		return
	}

	err = h.services.Ad.DeleteAd(ad.ID)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to delete the ad. Please try again.")
		bot.Send(msg)
		return
	}

	message := fmt.Sprintf("⚠️ Your ad titled '%s' has been deleted by an administrator.", ad.Title)

	h.NotifyUserAboutAdDeletion(bot, chatID, ad.SellerID, message)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ad with title '%s' deleted successfully.", ad.Title))
	bot.Send(msg)
}

func (h *AdminHandler) handleBroadcast(bot *tgbotapi.BotAPI, chatID int64, messageText string) {
	err := h.services.User.BroadcastMessage(messageText)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to broadcast message.")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Message broadcasted successfully!")
	bot.Send(msg)
}

func (h *AdminHandler) sendUnknownCommand(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Command not recognized.")
	bot.Send(msg)
}

func (h *AdminHandler) HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID

	switch {
	case data == "cancel":
		delete(h.userStates, chatID)
		h.sendAdminMainMenu(bot, chatID)
		return

	case strings.HasPrefix(data, "block_"):
		userID, err := strconv.Atoi(strings.TrimPrefix(data, "block_"))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Invalid user ID.")
			bot.Send(msg)
			return
		}
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
		userID, err := strconv.Atoi(strings.TrimPrefix(data, "grant_premium_"))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Invalid user ID.")
			bot.Send(msg)
			return
		}
		h.handleGrantPremium(bot, chatID, userID)

	case strings.HasPrefix(data, "delete_ad_"):
		userID, _ := strconv.Atoi(strings.TrimPrefix(data, "delete_ad_"))
		h.userStates[chatID] = fmt.Sprintf("waiting_for_ad_id_%s", userID)
		msg := tgbotapi.NewMessage(chatID, "Please enter the ID of the ad you want to delete:")
		bot.Send(msg)

	default:
		msg := tgbotapi.NewMessage(chatID, "Unknown command.")
		bot.Send(msg)
	}

	callbackResponse := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Request(callbackResponse)
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

func (h *AdminHandler) handleFindUserByUsername(bot *tgbotapi.BotAPI, chatID int64, username string) {
	user, err := h.services.User.GetUserByUsername(username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("User with username '%s' not found.", username))
		bot.Send(msg)
		return
	}

	userInfo := fmt.Sprintf(
		"User Info:\nID: %d\nName: %s\nBalance: %.2f\nRating: %.2f\nPremium: %t",
		user.TelegramID, user.Username, user.Balance, user.Rating, user.IsPremium,
	)

	msg := tgbotapi.NewMessage(chatID, userInfo)
	bot.Send(msg)
}

func (h *AdminHandler) NotifyUserAboutAdDeletion(bot *tgbotapi.BotAPI, chatID int64, sellerID int, messageText string) {
	err := h.services.User.BroadcastAboutDelete(sellerID, messageText)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to broadcast message.")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Message broadcasted successfully!")
	bot.Send(msg)
}
