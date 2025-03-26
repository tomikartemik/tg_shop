package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tg_shop/internal/model"
	"tg_shop/utils"
)

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID

	channelChatID := int64(-1002262695419)
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelChatID,
			UserID: telegramID,
		},
	})

	if err != nil || (member.Status != "member" && member.Status != "administrator" && member.Status != "creator") {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please subscribe to our channel to use the bot.")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Subscribe", fmt.Sprintf("https://t.me/%s", "+GtMFfelO1ko1ZWIy")),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("I am subscribed", "i_am_subscribed"),
			),
		)
		bot.Send(msg)
		return
	}

	// Остальная логика после проверки подписки...
	existingUser, err := h.services.GetUserInfoById(int(telegramID))
	if err == nil {
		user, err := h.services.GetUserById(int(telegramID))
		if err != nil {
			// Отправка видео
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("uploads/start.mp4"))
			video.Caption = "Welcome to Hell Market Bot!\n\nHell Market Bot is the place where you can safely purchase products from trusted sellers and list your own items for sale.\nOur goal is to make interaction between people as safe and fast as possible.\n\nEach listing is manually reviewed, ensuring 100% compliance and quality of the material you purchase.\n\nYou can learn more about how bot works by clicking on the article below this message. The guide will explain how this bot operates.\n\nAll important information and FAQ will be collected in the \"Important\" section in the main menu.\n\nDisclaimer: Our service works only with verified sellers. Any actions outside the law of any country will be stopped and condemned. All actions within this bot are conducted strictly within the bounds of the law."
			url := "https://telegra.ph/Instructions-for-working-with-the-bot-12-19"
			video.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("📘 Open Instructions", url),
				),
			)
			bot.Send(video)

			h.userStates[telegramID] = "username"

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your name:")
			bot.Send(msg)
			return
		}

		isBlocked := user.Banned

		if isBlocked {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are blocked and cannot use this bot.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Welcome back, %s!",
			existingUser.Username,
		))
		bot.Send(msg)
		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath("uploads/start.mp4"))
	video.Caption = "Welcome to Hell Market Bot!\n\nHell Market Bot is the place where you can safely purchase products from trusted sellers and list your own items for sale.\nOur goal is to make interaction between people as safe and fast as possible.\n\nEach listing is manually reviewed, ensuring 100% compliance and quality of the material you purchase.\n\nYou can learn more about how bot works by clicking on the article below this message. The guide will explain how this bot operates.\n\nAll important information and FAQ will be collected in the \"Important\" section in the main menu.\n\nDisclaimer: Our service works only with verified sellers. Any actions outside the law of any country will be stopped and condemned. All actions within this bot are conducted strictly within the bounds of the law."
	url := "https://telegra.ph/Instructions-for-working-with-the-bot-12-19"
	video.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📘 Open Instructions", url),
		),
	)
	bot.Send(video)

	h.userStates[telegramID] = "username"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your name:")
	bot.Send(msg)
}

func (h *Handler) HandleKeyboardButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, messageText string) {
	state, exists := h.userStates[update.Message.From.ID]
	if exists && strings.HasPrefix(state, "creating_ad") {
		h.handleAdCreation(bot, update, state, messageText)
		return
	}

	switch messageText {
	case "📝 Create Ad":
		h.userStates[update.Message.From.ID] = "creating_ad_title"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the title for your ad:")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("❌ Exit"),
			),
		)
		bot.Send(msg)
	case "👤 Profile":
		log.Printf("In switch: %s", messageText)

		user, err := h.services.GetUserById(int(update.Message.From.ID))
		if err != nil {
			log.Printf("Error fetching user profile: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error loading your profile.")
			bot.Send(msg)
			return
		}

		premiumStatus := "❌ Not Active"
		if user.IsPremium {
			premiumStatus = fmt.Sprintf("✅ Active until %s", user.ExpirePremium.Format("02 Jan 2006"))
		}

		escapedUsername := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, user.Username)
		escapedPremiumStatus := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, premiumStatus)

		profileMessage := fmt.Sprintf(
			"👤 *Your Profile:*\n"+
				"Id: `%d`\n"+
				"Name: `%s`\n"+
				"Balance: `%.2f$ (%.2f$)`\n"+
				"Rating: `%.2f`\n"+
				"Premium: `%s`",
			user.TelegramID, escapedUsername, user.Balance, user.HoldBalance, user.Rating, escapedPremiumStatus,
		)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📈 Add Balance", "add_balance"),
				tgbotapi.NewInlineKeyboardButtonData("📉 Request Payout", "request_payout"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✏️ Change Name", "change_name"),
				tgbotapi.NewInlineKeyboardButtonData("📄 My Ads", "my_ads"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🖼️ Change Photo", "change_photo"),
				tgbotapi.NewInlineKeyboardButtonData("📦My orders", "my_orders"),
			),
		)

		if user.PhotoURL == "" {
			log.Printf("User %d has no avatar. Sending profile without photo.", user.TelegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, profileMessage)
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		}

		photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(user.PhotoURL))
		photoMsg.Caption = profileMessage
		photoMsg.ParseMode = tgbotapi.ModeMarkdownV2
		photoMsg.ReplyMarkup = keyboard

		log.Printf("Sending profile with photo: %s", user.PhotoURL)
		if _, err := bot.Send(photoMsg); err != nil {
			log.Printf("Error sending profile photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to load your avatar. Here is your profile information.")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.Text = profileMessage
			bot.Send(msg)
		}

	case "💎 Premium":
		user, err := h.services.GetUserById(int(update.Message.From.ID))
		if err != nil {
			log.Printf("Error fetching user profile: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error loading your profile.")
			bot.Send(msg)
			return
		}
		msgText := ""
		if user.IsPremium {
			premiumStatus := fmt.Sprintf("✅ Active until %s", user.ExpirePremium.Format("02 Jan 2006"))
			msgText = fmt.Sprintf(
				"Your Premium: %s", premiumStatus,
			)
		} else {
			msgText = "Want to extend or purchase Premium? Contact the admin to get all the details and benefits!"
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Contact Admin", "https://t.me/Luc1ferTheDevil"),
				tgbotapi.NewInlineKeyboardButtonURL("❗️Terms of Premium", "https://telegra.ph/PREMIUM-02-20-3"),
			),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ReplyMarkup = keyboard

		bot.Send(msg)
	case "❗️Important":
		url := "https://telegra.ph/Instructions-for-working-with-the-bot-12-19"
		url_2 := "https://telegra.ph/Controversial-situations-Help-12-30"
		url_3 := "https://telegra.ph/User-Agreement-02-20-7"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Click the button below to view important information.")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("📘 Open Instructions", url),
				tgbotapi.NewInlineKeyboardButtonURL("📘 Controversial situations", url_2),
				tgbotapi.NewInlineKeyboardButtonURL("📘User Agreement", url_3),
			),
		)
		bot.Send(msg)
	case "🆘 Support":
		msgText := "If you have any questions or problems, our team is always ready to help you. You can contact the admin or the support to get the support you need. Also, if you have ideas or suggestions on how to improve our bot, we would love to hear them!\n\nPlease choose one of the following options:"

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Contact Admin", "https://t.me/Luc1ferTheDevil"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Contact Support", "https://t.me/hspquick"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Suggest an Idea", "https://t.me/Luc1ferTheDevil"),
			),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ReplyMarkup = keyboard

		bot.Send(msg)
	case "📄 Our channels":
		messageText := "Would be delighted if you check out our other projects listed below\\!\n\n" +
			"❗️All titles are clickable\\!\n\n" +
			"🔺 [HELL REFUND MAIN](https://t.me/\\+VtUPiZtDuX9hYTQy)\n\n" +
			"🔺 [HELL REFUND BACKUP](https://t.me/\\+ZOU4LSpBvwc5ZmRi)\n\n" +
			"🔺 [HELL REFUND CHAT](https://t.me/\\+3xhos0cIhTNhYmZi)\n\n" +
			"🔺 [HELL BOXING](https://t.me/\\+X9-Ql8LQVDYyYmI6)\n\n" +
			"🔺 [HELL CHECKIP BOT](https://t.me/hellcheckip_bot)"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		msg.ParseMode = "MarkdownV2"

		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Command not recognized.")
		bot.Send(msg)
	}
}

func (h *Handler) HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	//
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID
	messageID := callbackQuery.Message.MessageID

	if strings.HasPrefix(data, "approve_ad_") {
		parts := strings.Split(data, "_")
		adID, _ := strconv.Atoi(parts[2])
		groupID, _ := strconv.ParseInt(parts[3], 10, 64)

		if err := h.services.Ad.ApproveAd(adID); err != nil {
			log.Printf("Failed to approve ad: %s", err)
			return
		}

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			groupID,
			messageID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		if _, err := bot.Send(editMarkup); err != nil {
			log.Printf("Failed to remove buttons: %v", err)
		}

		ad, err := h.services.Ad.GetAdByIDTg(adID)
		if err == nil {
			h.NotifyUser(bot, ad.SellerID, ad, true)
		}

	} else if strings.HasPrefix(data, "reject_ad_") {
		parts := strings.Split(data, "_")
		adID, _ := strconv.Atoi(parts[2])
		groupID, _ := strconv.ParseInt(parts[3], 10, 64)

		if err := h.services.Ad.RejectAd(adID); err != nil {
			log.Printf("Failed to reject ad: %s", err)
			return
		}

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			groupID,
			messageID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		if _, err := bot.Send(editMarkup); err != nil {
			log.Printf("Failed to remove buttons: %v", err)
		}

		ad, err := h.services.Ad.GetAdByIDTg(adID)
		if err == nil {
			h.NotifyUser(bot, ad.SellerID, ad, false)
		}

	} else if strings.HasPrefix(data, "approve_payout_") {
		parts := strings.Split(data, "_")
		log.Printf("Parts: %v", parts)
		payoutID, _ := strconv.Atoi(parts[2])
		groupID, _ := strconv.ParseInt(parts[3], 10, 64)

		payout, err := h.services.Payout.GetPayoutByID(payoutID)
		if err != nil {
			log.Printf("Error fetching payout for payoutID %d: %v", payoutID, err)
			return
		}

		user, err := h.services.GetUserById(payout.TelegramID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			return
		}

		newBalance := user.Balance - payout.Amount
		if newBalance < 0 {
			log.Printf("Insufficient balance for payout request ID %d", payoutID)
			msg := tgbotapi.NewMessage(chatID, "❌ Error: Insufficient balance to process the payout.")
			bot.Send(msg)
			return
		}

		err = h.services.ChangeBalance(user.TelegramID, newBalance)
		if err != nil {
			log.Printf("Error updating user balance: %v", err)
			msg := tgbotapi.NewMessage(chatID, "❌ Error: Failed to update user balance.")
			bot.Send(msg)
			return
		}

		err = h.services.Payout.ApprovePayoutRequest(payoutID)
		if err != nil {
			log.Printf("Error approving payout: %v", err)
			msg := tgbotapi.NewMessage(chatID, "❌ Error: Failed to approve payout.")
			bot.Send(msg)
			return
		}

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			groupID,
			messageID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		if _, err := bot.Send(editMarkup); err != nil {
			log.Printf("Failed to remove buttons: %v", err)
		}

		h.NotifyPayout(bot, user, payout.Amount, true)

	} else if strings.HasPrefix(data, "reject_payout_") {
		parts := strings.Split(data, "_")
		payoutID, _ := strconv.Atoi(parts[2])
		groupID, _ := strconv.ParseInt(parts[3], 10, 64)

		payout, err := h.services.Payout.GetPayoutByID(payoutID)
		if err != nil {
			log.Printf("Error fetching payout for payout ID %d: %v", payoutID, err)
			return
		}

		err = h.services.Payout.RejectPayoutRequest(payoutID)
		if err != nil {
			log.Printf("Error rejecting payout: %v", err)
			return
		}

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			groupID,
			messageID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		if _, err := bot.Send(editMarkup); err != nil {
			log.Printf("Failed to remove buttons: %v", err)
		}

		user, err := h.services.GetUserById(payout.TelegramID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			return
		}
		h.NotifyPayout(bot, user, payout.Amount, false)

	} else if strings.HasPrefix(data, "rate_") {
		parts := strings.Split(data, "_")
		if len(parts) != 3 {
			log.Printf("Invalid rating callback data: %s", data)
			return
		}

		sellerID, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Failed to parse sellerID: %v", err)
			return
		}

		rating, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("Failed to parse rating: %v", err)
			return
		}

		err = h.services.User.ChangeRating(sellerID, rating)
		if err != nil {
			log.Printf("Failed to change seller rating: %v", err)
			msg := tgbotapi.NewMessage(chatID, "An error occurred when changing the seller's rating.")
			bot.Send(msg)
			return
		}

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			messageID,
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			},
		)
		if _, err := bot.Send(editMarkup); err != nil {
			log.Printf("Failed to remove buttons: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Thanks! You rated the seller on %d.", rating))
		bot.Send(msg)
		return
	} else {
		switch data {
		case "i_am_subscribed":
			telegramID := callbackQuery.From.ID

			channelChatID := int64(-1002262695419)
			member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					ChatID: channelChatID,
					UserID: telegramID,
				},
			})

			if err != nil || (member.Status != "member" && member.Status != "administrator" && member.Status != "creator") {
				msg := tgbotapi.NewMessage(chatID, "You are not subscribed yet. Please subscribe to the channel.")
				bot.Send(msg)
				return
			}

			update := tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: chatID,
					},
					From: &tgbotapi.User{
						ID: telegramID,
					},
					Text: "/start",
				},
			}
			h.HandleStart(bot, update) // Повторный вызов команды /start
			return
		case "add_balance":
			h.userStates[callbackQuery.From.ID] = "adding_balance"

			msg := tgbotapi.NewMessage(chatID, "Enter the amount to top up or press 'Cancel':")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Cancel"),
				),
			)
			bot.Send(msg)

		case "request_payout":
			user, err := h.services.GetUserById(int(callbackQuery.From.ID))
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				msg := tgbotapi.NewMessage(chatID, "Failed to load your profile. Please try again later.")
				bot.Send(msg)
				return
			}

			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
				"Your balance: %.2f$\nMinimum withdrawal amount: $50\nEnter the amount you want to withdraw or press 'Cancel':",
				user.Balance,
			))
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Cancel"),
				),
			)
			bot.Send(msg)

			// Устанавливаем состояние для запроса суммы
			h.userStates[callbackQuery.From.ID] = "requesting_payout_amount"

		case "change_name":
			msg := tgbotapi.NewMessage(chatID, "Please enter your new name:")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Cancel"),
				),
			)
			bot.Send(msg)
			h.userStates[callbackQuery.From.ID] = "changing_name"

		case "my_ads":
			ads, err := h.services.Ad.GetAdsByUserID(int(callbackQuery.From.ID))
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Error loading your ads. Please try again later.")
				bot.Send(msg)
				return
			}

			if len(ads) == 0 {
				msg := tgbotapi.NewMessage(chatID, "You have no ads.")
				bot.Send(msg)
				return
			}

			adsMessage := "📄 *Your Ads:*\n"
			for _, ad := range ads {
				adsMessage += fmt.Sprintf(
					"ID: %d\nTitle: %s\nPrice: %.2f$\nStock: %d\nStatus: %s\n\n",
					ad.ID, ad.Title, ad.Price, ad.Stock, ad.Status,
				)
			}

			deleteButton := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("❌ Delete Ad", "delete_ad"),
				),
			)

			msg := tgbotapi.NewMessage(chatID, adsMessage)
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = deleteButton
			bot.Send(msg)

		case "delete_ad":
			h.userStates[callbackQuery.From.ID] = "deleting_ad"

			msg := tgbotapi.NewMessage(chatID, "Enter the ID of the ad you want to delete or press 'Cancel':")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Cancel"),
				),
			)
			bot.Send(msg)

		case "my_orders":
			user, err := h.services.GetUserById(int(callbackQuery.From.ID))
			if err != nil {
				log.Printf("Error fetching user orders: %v", err)
				msg := tgbotapi.NewMessage(chatID, "Error loading your orders. Please try again later.")
				bot.Send(msg)
				return
			}

			if len(user.Purchased) == 0 {
				msg := tgbotapi.NewMessage(chatID, "You have no purchases.")
				bot.Send(msg)
				return
			}

			ordersMessage := "🛒 *Your Orders:*\n"
			for _, ad := range user.Purchased {
				ordersMessage += fmt.Sprintf(
					"\n*Title:* %s\n*Price:* %.2f\n*Description:* %s\n\n",
					ad.Title, ad.Price, ad.Description,
				)
			}

			msg := tgbotapi.NewMessage(chatID, ordersMessage)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		case "change_photo":
			msg := tgbotapi.NewMessage(chatID, "Please upload your new profile picture:")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("❌ Cancel"),
				),
			)
			bot.Send(msg)
			h.userStates[callbackQuery.From.ID] = "changing_photo"
		default:
			user, err := h.services.GetUserById(int(chatID))
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				msg := tgbotapi.NewMessage(chatID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}

			isBlocked := user.Banned

			if isBlocked {
				msg := tgbotapi.NewMessage(chatID, "You are blocked and cannot use this bot.")
				bot.Send(msg)
				return
			}
		}
	}
}

func (h *Handler) HandleUserInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	messageText := strings.TrimSpace(update.Message.Text)

	if messageText == "❌ Cancel" || messageText == "❌ Exit" {
		delete(h.userStates, telegramID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Operation has been canceled.")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		bot.Send(msg)
		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	log.Printf("User %d state: %s", telegramID, h.userStates[telegramID])
	log.Printf("Received message: %s", messageText)

	if h.userStates[telegramID] == "username" {
		username := messageText

		user := model.User{
			TelegramID: int(telegramID),
			Username:   username,
		}

		_, err := h.services.CreateOrUpdateUser(user)
		if err != nil {
			log.Printf("Error updating username: %v", err)
			if strings.Contains(err.Error(), "duplicate") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Username is already taken.")
				bot.Send(msg)
			}
			return
		}

		h.userStates[telegramID] = "uploading_avatar"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your name is saved. Please upload a profile picture or press 'Skip':")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("✅ Skip"),
			),
		)
		bot.Send(msg)
		return
	} else if h.userStates[telegramID] == "deleting_ad" {
		adID, err := strconv.Atoi(messageText)
		if err != nil || adID <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid ad ID. Please enter a valid numeric ID:")
			bot.Send(msg)
			return
		}

		ad, err := h.services.Ad.GetAdByIDTg(adID)
		if err != nil || ad.SellerID != int(telegramID) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This ad does not belong to you or does not exist.")
			bot.Send(msg)
			return
		}

		err = h.services.Ad.DeleteAd(adID)
		if err != nil {
			log.Printf("Error deleting ad: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to delete the ad. Please try again later.")
			bot.Send(msg)
			return
		}

		delete(h.userStates, telegramID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ad with ID %d has been successfully deleted.", adID))
		bot.Send(msg)
		h.sendMainMenu(bot, update.Message.Chat.ID)
	} else if strings.HasPrefix(h.userStates[telegramID], "requesting_payout") {
		// Если пользователь нажал "Отмена"
		if strings.TrimSpace(messageText) == "❌ Cancel" {
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Payout operation has been canceled.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}

		// Обработка в зависимости от текущего состояния
		if h.userStates[telegramID] == "requesting_payout_amount" {
			amount, err := strconv.ParseFloat(messageText, 64)
			if err != nil || amount <= 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid amount. Please enter a positive number.")
				bot.Send(msg)
				return
			}

			// Проверка на минимальную сумму вывода
			if amount < 50 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Minimum withdrawal amount is $50.")
				bot.Send(msg)
				return
			}

			user, err := h.services.GetUserById(int(telegramID))
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to load your profile. Please try again later.")
				bot.Send(msg)
				return
			}

			if amount > user.Balance {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Insufficient balance.")
				bot.Send(msg)
				return
			}

			// Сохраняем сумму в состоянии
			h.userStates[telegramID] = fmt.Sprintf("requesting_payout_username:%f", amount)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your username (must start with @):")
			bot.Send(msg)

		} else if strings.HasPrefix(h.userStates[telegramID], "requesting_payout_username:") {
			if !strings.HasPrefix(messageText, "@") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid username. It must start with @. Please try again:")
				bot.Send(msg)
				return
			}

			// Извлекаем сумму из состояния
			parts := strings.Split(h.userStates[telegramID], ":")
			if len(parts) != 2 {
				log.Printf("Invalid state format: %s", h.userStates[telegramID])
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}

			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Printf("Failed to parse amount: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}

			// Сохраняем username в состоянии
			h.userStates[telegramID] = fmt.Sprintf("requesting_payout_wallet:%f:%s", amount, messageText)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your USDT TRC20 wallet address:")
			bot.Send(msg)

		} else if strings.HasPrefix(h.userStates[telegramID], "requesting_payout_wallet:") {
			// Извлекаем сумму и username из состояния
			parts := strings.Split(h.userStates[telegramID], ":")
			if len(parts) != 3 {
				log.Printf("Invalid state format: %s", h.userStates[telegramID])
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}

			amount, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Printf("Failed to parse amount: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}

			username := parts[2]
			wallet := messageText

			// Создаем запрос на вывод
			payoutID, err := h.services.Payout.CreatePayoutRequest(int(telegramID), amount)
			if err != nil {
				log.Printf("Error creating payout request: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to create payout request. Please try again later.")
				bot.Send(msg)
				return
			}

			log.Printf("Payout request created with ID: %d", payoutID)

			// Отправка сообщения в группу модерации
			user, err := h.services.GetUserById(int(telegramID))
			if err != nil {
				log.Printf("Error fetching user: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
				bot.Send(msg)
				return
			}
			payoutGroupID, _ := strconv.ParseInt(os.Getenv("GROUP_WITHDRAWAL_ID"), 10, 64)
			messageID, err := h.SendPayoutRequestToModeration(bot, user, amount, payoutGroupID, payoutID, username, wallet)
			if err != nil {
				log.Printf("Error sending payout request to moderation group: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to notify moderators. Please try again later.")
				bot.Send(msg)
				return
			}

			_ = messageID

			delete(h.userStates, telegramID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your payout request has been submitted for moderation.")
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}
	} else if h.userStates[telegramID] == "changing_name" {
		newName := messageText

		if strings.TrimSpace(messageText) == "❌ Cancel" {
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Change name operation has been canceled.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}

		if len(newName) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Name cannot be empty. Please enter a valid name:")
			bot.Send(msg)
			return
		}

		updatedUser := model.User{
			TelegramID: int(telegramID),
			Username:   newName,
		}

		_, err := h.services.CreateOrUpdateUser(updatedUser)
		if err != nil {
			log.Printf("Error updating username: %v", err)
			if strings.Contains(err.Error(), "duplicate") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Username is already taken.")
				bot.Send(msg)
			}
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your name has been updated to: %s", newName))
		bot.Send(msg)

		delete(h.userStates, telegramID)
		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	} else if h.userStates[telegramID] == "uploading_avatar" {
		// Если пользователь нажал "Skip"
		if strings.TrimSpace(messageText) == "✅ Skip" {
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Profile picture upload skipped. Your profile has been created.")
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}

		if update.Message.Photo == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a valid photo or press 'Skip'.")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("✅ Skip"),
				),
			)
			bot.Send(msg)
			return
		}

		photo := update.Message.Photo[len(update.Message.Photo)-1]

		fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
		file, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Printf("Failed to get file from Telegram: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}

		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)

		response, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to download file: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}
		defer response.Body.Close()

		fileData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Failed to read file data: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}

		fileName := fmt.Sprintf("%d_avatar.jpg", update.Message.From.ID)
		filePath, err := utils.SaveFile(fileData, fileName, "./uploads")
		if err != nil {
			log.Printf("Error saving photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save the photo. Try again.")
			bot.Send(msg)
			return
		}

		updatedUser := model.User{
			TelegramID: int(update.Message.From.ID),
			PhotoURL:   filePath,
		}

		_, err = h.services.CreateOrUpdateUser(updatedUser)
		if err != nil {
			log.Printf("Error updating user avatar: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save the avatar. Please try again.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your profile picture is saved!")
		bot.Send(msg)

		delete(h.userStates, update.Message.From.ID)
		h.sendMainMenu(bot, update.Message.Chat.ID)
	} else if h.userStates[telegramID] == "changing_photo" {

		if strings.TrimSpace(messageText) == "❌ Cancel" {
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Change photo operation has been canceled.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}

		if update.Message.Photo == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a valid photo.")
			bot.Send(msg)
			return
		}

		photo := update.Message.Photo[len(update.Message.Photo)-1]

		fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
		file, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Printf("Failed to get file from Telegram: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}

		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)

		response, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to download file: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}
		defer response.Body.Close()

		fileData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Failed to read file data: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			bot.Send(msg)
			return
		}

		fileName := fmt.Sprintf("%d_avatar.jpg", update.Message.From.ID)
		filePath, err := utils.SaveFile(fileData, fileName, "./uploads")
		if err != nil {
			log.Printf("Error saving photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save the photo. Try again.")
			bot.Send(msg)
			return
		}

		updatedUser := model.User{
			TelegramID: int(update.Message.From.ID),
			PhotoURL:   filePath,
		}

		_, err = h.services.CreateOrUpdateUser(updatedUser)
		if err != nil {
			log.Printf("Error updating user photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to update your profile picture. Please try again.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your profile picture has been updated successfully!")
		bot.Send(msg)

		delete(h.userStates, update.Message.From.ID)
		h.sendMainMenu(bot, update.Message.Chat.ID)
	} else if h.userStates[telegramID] == "adding_balance" {
		// If the user pressed "Cancel"
		if strings.TrimSpace(messageText) == "❌ Cancel" {
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Top-up operation has been canceled.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}

		// Validate the entered amount
		amount, err := strconv.ParseFloat(messageText, 64)
		if err != nil || amount <= 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid amount. Please enter a number greater than 0.")
			bot.Send(msg)
			return
		}

		// Create an invoice through CryptoCloud
		link, err := h.services.CryptoCloud.CreateInvoice(amount, int(telegramID))
		if err != nil {
			log.Printf("Error creating invoice: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred while creating the invoice. Please try again.")
			bot.Send(msg)
			return
		}

		// Clear user state
		delete(h.userStates, telegramID)

		// Create an inline button with the payment link
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("💳 Pay Now", link),
			),
		)

		// Send the payment button to the user
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("An invoice for %.2f has been created. Click the button below to complete the payment.", amount))
		msg.ReplyMarkup = keyboard
		bot.Send(msg)

		h.sendMainMenu(bot, update.Message.Chat.ID)
	} else {
		user, err := h.services.GetUserById(int(telegramID))
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
			bot.Send(msg)
			return
		}

		isBlocked := user.Banned

		if isBlocked {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are blocked and cannot use this bot.")
			bot.Send(msg)
			return
		}
		h.HandleKeyboardButton(bot, update, messageText)
	}
}

func (h *Handler) sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	menuMessage := "Choose an action from the menu below:"

	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝 Create Ad"),
			tgbotapi.NewKeyboardButton("👤 Profile"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💎 Premium"),
			tgbotapi.NewKeyboardButton("❗️Important"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🆘 Support"),
			tgbotapi.NewKeyboardButton("📄 Our channels"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, menuMessage)
	msg.ReplyMarkup = menuKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending main menu: %v", err)
	}
}

func (h *Handler) handleAdCreation(bot *tgbotapi.BotAPI, update tgbotapi.Update, state, messageText string) {
	telegramID := update.Message.From.ID

	if messageText == "❌ Exit" {
		delete(h.tempAdData, telegramID)
		delete(h.userStates, telegramID)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You have exited the ad creation process.")
		bot.Send(msg)

		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	switch state {
	case "creating_ad_title":
		if len(messageText) > 100 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The title is too long. Please enter a title with a maximum of 100 characters.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		h.tempAdData[telegramID] = model.Ad{Title: messageText}
		h.userStates[telegramID] = "creating_ad_description"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the description for your ad:")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_description":
		if len(messageText) > 700 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The description is too long. Please enter a description with a maximum of 700 characters.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		ad := h.tempAdData[telegramID]
		paragraphs := strings.Split(messageText, "\n")
		formattedDescription := strings.Join(paragraphs, "\n")
		ad.Description = formattedDescription
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_price"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the price for your ad (Only Dollars $):")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_price":
		price, err := strconv.ParseFloat(messageText, 64)
		if err != nil || price < 0 { // Добавлена проверка price <= 0
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid price. Please enter a positive numeric value:")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}
		ad := h.tempAdData[telegramID]
		ad.Price = price
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_stock"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the stock quantity for your ad:")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("♾️ Unlimited"),
				tgbotapi.NewKeyboardButton("❌ Exit"),
			),
		)
		bot.Send(msg)

	case "creating_ad_stock":
		if messageText == "♾️ Unlimited" {
			ad := h.tempAdData[telegramID]
			ad.Stock = 9999999
			h.tempAdData[telegramID] = ad
			h.userStates[telegramID] = "creating_ad_category"

			categories, err := h.services.Category.GetCategoryList()
			if err != nil {
				log.Printf("Error fetching categories: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to load categories. Please try again later.")
				bot.Send(msg)
				return
			}

			if len(categories) == 0 {
				log.Println("No categories found.")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No categories available at the moment.")
				bot.Send(msg)
				return
			}

			var categoryList strings.Builder
			categoryList.WriteString("📋 *Available Categories:*\n")
			for _, category := range categories {
				categoryList.WriteString(fmt.Sprintf("%d - %s\n", category.ID, category.Name))
			}
			categoryList.WriteString("\nPlease enter the ID of the category you want to choose:")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, categoryList.String())
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		stock, err := strconv.Atoi(messageText)
		if err != nil || stock <= 0 { // Добавлена проверка stock <= 0
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Invalid quantity. Please enter a positive number or select '♾️ Unlimited':")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("♾️ Unlimited"),
					tgbotapi.NewKeyboardButton("❌ Exit"),
				),
			)
			bot.Send(msg)
			return
		}

		ad := h.tempAdData[telegramID]
		ad.Stock = stock
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_category"

		categories, err := h.services.Category.GetCategoryList()
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to load categories. Please try again later.")
			bot.Send(msg)
			return
		}

		if len(categories) == 0 {
			log.Println("No categories found.")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No categories available at the moment.")
			bot.Send(msg)
			return
		}

		var categoryList strings.Builder
		categoryList.WriteString("📋 *Available Categories:*\n")
		for _, category := range categories {
			categoryList.WriteString(fmt.Sprintf("%d - %s\n", category.ID, category.Name))
		}
		categoryList.WriteString("\nPlease enter the ID of the category you want to choose:")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, categoryList.String())
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_category":
		categoryID, err := strconv.Atoi(messageText)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid category ID. Please enter a numeric value:")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}
		ad := h.tempAdData[telegramID]
		ad.CategoryID = categoryID
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_photo"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a photo for your ad:")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_photo":
		if update.Message.Photo == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a valid photo.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		photo := update.Message.Photo[len(update.Message.Photo)-1]

		fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
		file, err := bot.GetFile(fileConfig)
		if err != nil {
			log.Printf("Failed to get file from Telegram: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		response, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to download file: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}
		defer response.Body.Close()

		fileData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Failed to read file data: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process the photo. Try again.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		filePath, err := utils.SaveFile(fileData, "ad_photo.jpg", "./uploads")
		if err != nil {
			log.Printf("Error saving photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save the photo. Try again.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		ad := h.tempAdData[update.Message.From.ID]
		ad.PhotoURL = filePath
		h.tempAdData[update.Message.From.ID] = ad

		h.userStates[update.Message.From.ID] = "creating_ad_files"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload any additional files for your ad.")
		msg.ReplyMarkup = getAdCreationButtons("creating_ad_files")
		bot.Send(msg)

	case "creating_ad_files":
		ad := h.tempAdData[telegramID]

		if messageText == "✅ Skip" {
			h.userStates[telegramID] = "creating_ad_finish"
		} else if update.Message.Document != nil {
			file := update.Message.Document
			fileData, err := downloadFileFromTg(bot, file.FileID)
			if err != nil {
				log.Fatalf("Error downloading file: %v", err)
			}

			filePath, err := utils.SaveFile(fileData, file.FileName, "./uploads")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to upload file. Please try again.")
				msg.ReplyMarkup = getExitKeyboard()
				bot.Send(msg)
				return
			}
			ad.Files = filePath
			h.userStates[telegramID] = "creating_ad_finish"
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a valid file.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		h.tempAdData[telegramID] = ad

		category, err := h.services.Category.GetCategoryById(ad.CategoryID)

		if err != nil {
			log.Printf("Error fetching category: %v", err)
		}

		photoTxt := fmt.Sprintf(
			"📢 *Ad for Moderation:*\n"+
				"**Title:** %s\n"+
				"**Description:** %s\n"+
				"**Price:** %.2f$\n"+
				"**Stock:** %d\n"+
				"**Category:** %s",
			ad.Title, ad.Description, ad.Price, ad.Stock, category.Name,
		)
		photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(ad.PhotoURL))
		photo.Caption = photoTxt
		photo.ParseMode = "Markdown"
		messageTxt := "Do you want to submit this ad? Use the buttons below:"
		message := tgbotapi.NewMessage(update.Message.Chat.ID, messageTxt)
		message.ReplyMarkup = getAdCreationButtons("creating_ad_finish")

		bot.Send(photo)

		bot.Send(message)

	case "creating_ad_finish":
		if messageText == "❌ Cancel" {
			delete(h.tempAdData, telegramID)
			delete(h.userStates, telegramID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ad creation canceled.")
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
			return
		}
		if messageText == "✅ Confirm" {
			ad := h.tempAdData[telegramID]
			ad.SellerID = int(telegramID)
			createdAd, err := h.services.Ad.CreateAd(ad)
			if err != nil {
				log.Printf("Error creating ad: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error creating ad. Please try again later.")
				bot.Send(msg)
				h.sendMainMenu(bot, update.Message.Chat.ID)
				return
			}
			delete(h.tempAdData, telegramID)
			delete(h.userStates, telegramID)
			moderationGroupID, _ := strconv.ParseInt(os.Getenv("GROUP_MODERATION_ID"), 10, 64)
			ourGroupID, _ := strconv.ParseInt(os.Getenv("OUR_GROUP"), 10, 64)
			messageID, err := h.SendAdToModeration(bot, createdAd, moderationGroupID)
			messageID2, err2 := h.SendAdToOurGroup(bot, createdAd, ourGroupID)
			if err != nil {
				log.Printf("Error sending ad to moderation group: %v", err)
			}
			if err2 != nil {
				log.Printf("Error sending ad to our group: %v", err)
			}

			_ = messageID
			_ = messageID2
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
				"Your ad '%s' has been submitted for moderation. Ad ID: %d",
				createdAd.Title, createdAd.ID,
			))
			bot.Send(msg)
			h.sendMainMenu(bot, update.Message.Chat.ID)
		}
	}
}

func getAdCreationButtons(state string) tgbotapi.ReplyKeyboardMarkup {
	switch state {
	case "creating_ad_files":
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("❌ Exit"),
			),
		)
	case "creating_ad_finish":
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("✅ Confirm"),
				tgbotapi.NewKeyboardButton("❌ Cancel"),
			),
		)
	default:
		return tgbotapi.NewReplyKeyboard()
	}
}

func getExitKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Exit"),
		),
	)
}

func downloadFileFromTg(bot *tgbotapi.BotAPI, fileID string) ([]byte, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %v", err)
	}

	return fileData, nil
}

func (h *Handler) SendAdToModeration(bot *tgbotapi.BotAPI, ad model.Ad, moderationGroupID int64) (int, error) {
	category, err := h.services.Category.GetCategoryById(ad.CategoryID)
	if err != nil {
		log.Printf("Error fetching category: %v", err)
		return 0, err
	}

	messageText := fmt.Sprintf(
		"📢 *Ad for Moderation:*\n"+
			"**Title:** %s\n"+
			"**Description:** %s\n"+
			"**Price:** %.2f$\n"+
			"**Stock:** %d\n"+
			"**Category:** %s\n"+
			"**Seller:** %d\n",
		ad.Title, ad.Description, ad.Price, ad.Stock, category.Name, ad.SellerID,
	)

	// Отправляем фото
	photo := tgbotapi.NewPhoto(moderationGroupID, tgbotapi.FilePath(ad.PhotoURL))
	photo.Caption = messageText
	photo.ParseMode = "Markdown"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Approve", fmt.Sprintf("approve_ad_%d_%d", ad.ID, moderationGroupID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Reject", fmt.Sprintf("reject_ad_%d_%d", ad.ID, moderationGroupID)),
		),
	)

	photo.ReplyMarkup = keyboard
	sentMsg, err := bot.Send(photo)
	if err != nil {
		return 0, err
	}

	if ad.Files != "" {
		file := tgbotapi.NewDocument(moderationGroupID, tgbotapi.FilePath(ad.Files))
		_, err := bot.Send(file)
		if err != nil {
			log.Printf("Failed to send file to moderation group: %v", err)
			return sentMsg.MessageID, err
		}
	}

	return sentMsg.MessageID, nil
}

func (h *Handler) SendAdToOurGroup(bot *tgbotapi.BotAPI, ad model.Ad, moderationGroupID int64) (int, error) {
	category, err := h.services.Category.GetCategoryById(ad.CategoryID)
	if err != nil {
		log.Printf("Error fetching category: %v", err)
		return 0, err
	}

	messageText := fmt.Sprintf(
		"📢 *Ad for Moderation:*\n"+
			"**Title:** %s\n"+
			"**Description:** %s\n"+
			"**Price:** %.2f$\n"+
			"**Stock:** %d\n"+
			"**Category:** %s\n"+
			"**Seller:** %d\n",
		ad.Title, ad.Description, ad.Price, ad.Stock, category.Name, ad.SellerID,
	)

	// Отправляем фото
	photo := tgbotapi.NewPhoto(moderationGroupID, tgbotapi.FilePath(ad.PhotoURL))
	photo.Caption = messageText
	photo.ParseMode = "Markdown"

	sentMsg, err := bot.Send(photo)
	if err != nil {
		return 0, err
	}

	if ad.Files != "" {
		file := tgbotapi.NewDocument(moderationGroupID, tgbotapi.FilePath(ad.Files))
		_, err := bot.Send(file)
		if err != nil {
			log.Printf("Failed to send file to moderation group: %v", err)
			return sentMsg.MessageID, err
		}
	}

	return sentMsg.MessageID, nil
}

func (h *Handler) SendPayoutRequestToModeration(bot *tgbotapi.BotAPI, user model.User, amount float64, payoutGroupID int64, payoutID int, username, wallet string) (int, error) {
	messageText := fmt.Sprintf(
		"💸 *Payout Request:*\n"+
			"**User:** %s\n"+
			"**Telegram ID:** %d\n"+
			"**Username:** `%s`\n"+
			"**USDT TRC20 Wallet:** `%s`\n"+
			"**Amount:** %.2f$\n",
		user.Username, user.TelegramID, username, wallet, amount,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Approve", fmt.Sprintf("approve_payout_%d_%d", payoutID, payoutGroupID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Reject", fmt.Sprintf("reject_payout_%d_%d", payoutID, payoutGroupID)),
		),
	)

	msg := tgbotapi.NewMessage(payoutGroupID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	sentMsg, err := bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return sentMsg.MessageID, nil
}
func (h *Handler) NotifyUser(bot *tgbotapi.BotAPI, userID int, ad model.Ad, approved bool) {
	var messageText string
	if approved {
		messageText = fmt.Sprintf("🎉 Your ad '%s' has been approved and is now visible!", ad.Title)
	} else {
		messageText = fmt.Sprintf("❌ Your ad '%s' has been rejected.", ad.Title)
	}

	bot.Send(tgbotapi.NewMessage(int64(userID), messageText))
}

func (h *Handler) NotifyPayout(bot *tgbotapi.BotAPI, user model.User, amount float64, status bool) {
	var messageText string

	if status {
		// Уведомление об успешной выплате
		messageText = fmt.Sprintf(
			"🎉 Your payout request for %.2f$ has been successfully processed! The funds should arrive shortly.",
			amount,
		)
	} else {
		// Уведомление об отклонении выплаты
		messageText = fmt.Sprintf(
			"❌ Your payout request for %.2f$ has been declined. Please contact support for more details.",
			amount,
		)
	}

	// Отправка уведомления пользователю
	msg := tgbotapi.NewMessage(int64(user.TelegramID), messageText)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send payout notification to user %d: %v", user.TelegramID, err)
	}
}

func (h *Handler) SendRatingRequest(bot *tgbotapi.BotAPI, chatID int64, sellerID int) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Please, rate seller (ID: %d):", sellerID))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", fmt.Sprintf("rate_%d_1", sellerID)),
			tgbotapi.NewInlineKeyboardButtonData("2", fmt.Sprintf("rate_%d_2", sellerID)),
			tgbotapi.NewInlineKeyboardButtonData("3", fmt.Sprintf("rate_%d_3", sellerID)),
			tgbotapi.NewInlineKeyboardButtonData("4", fmt.Sprintf("rate_%d_4", sellerID)),
			tgbotapi.NewInlineKeyboardButtonData("5", fmt.Sprintf("rate_%d_5", sellerID)),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send rating request: %v", err)
	}
}
