package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tg_shop/internal/model"
	"tg_shop/utils"
)

func (h *Handler) HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID

	existingUser, err := h.services.GetUserById(int(telegramID))
	if err == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Welcome back, %s!",
			existingUser.Username,
		))
		bot.Send(msg)
		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	}

	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Printf("Error checking user existence: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred. Please try again later.")
		bot.Send(msg)
		return
	}

	h.userStates[telegramID] = "username"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your name:")
	bot.Send(msg)
}

func (h *Handler) HandleUserInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	telegramID := update.Message.From.ID
	messageText := strings.TrimSpace(update.Message.Text)

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
			log.Printf("Error creating/updating user: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error saving your name. Please try again.")
			bot.Send(msg)
			return
		}

		h.userStates[telegramID] = "uploading_avatar"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your name is saved. Please upload a profile picture:")
		bot.Send(msg)
		return
	} else if h.userStates[telegramID] == "changing_name" {
		newName := messageText

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to update your name. Please try again later.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your name has been updated to: %s", newName))
		bot.Send(msg)

		delete(h.userStates, telegramID)
		h.sendMainMenu(bot, update.Message.Chat.ID)
		return
	} else if h.userStates[telegramID] == "uploading_avatar" {
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
	} else {
		h.HandleKeyboardButton(bot, update, messageText)
	}
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

		profileMessage := fmt.Sprintf(
			"👤 *Your Profile:*\n"+
				"Name: %s\n"+
				"Balance: %.2f\n"+
				"Rating: %.2f\n"+
				"Premium: %s",
			user.Username, user.Balance, user.Rating, premiumStatus,
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
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		}

		photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(user.PhotoURL))
		photoMsg.Caption = profileMessage
		photoMsg.ParseMode = "Markdown"
		photoMsg.ReplyMarkup = keyboard

		log.Printf("Sending profile with photo: %s", user.PhotoURL)
		if _, err := bot.Send(photoMsg); err != nil {
			log.Printf("Error sending profile photo: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to load your avatar. Here is your profile information.")
			msg.ParseMode = "Markdown"
			msg.Text = profileMessage
			bot.Send(msg)
		}

	case "💎 Premium":
		context.TODO()
	case "❗️Important":
		url := "https://telegra.ph/Instructions-for-working-with-the-bot-12-14"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Click the button below to view important information.")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("📘 Open Instructions", url),
			),
		)
		bot.Send(msg)
	case "🆘 Support":
		context.TODO()
	case "📄 Our channels":
		context.TODO()
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Command not recognized.")
		bot.Send(msg)
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
		h.tempAdData[telegramID] = model.Ad{Title: messageText}
		h.userStates[telegramID] = "creating_ad_description"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the description for your ad:")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_description":
		ad := h.tempAdData[telegramID]
		ad.Description = messageText
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_price"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the price for your ad:")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_price":
		price, err := strconv.ParseFloat(messageText, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid price. Please enter a numeric value:")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}
		ad := h.tempAdData[telegramID]
		ad.Price = price
		h.tempAdData[telegramID] = ad
		h.userStates[telegramID] = "creating_ad_stock"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the stock quantity for your ad:")
		msg.ReplyMarkup = getExitKeyboard()
		bot.Send(msg)

	case "creating_ad_stock":
		stock, err := strconv.Atoi(messageText)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid stock. Please enter a numeric value:")
			msg.ReplyMarkup = getExitKeyboard()
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload any additional files for your ad or click 'Skip'.")
		msg.ReplyMarkup = getAdCreationButtons("creating_ad_files")
		bot.Send(msg)

	case "creating_ad_files":
		ad := h.tempAdData[telegramID]

		if messageText == "✅ Skip" {
			h.userStates[telegramID] = "creating_ad_finish"
		} else if update.Message.Document != nil {
			file := update.Message.Document
			filePath, err := utils.SaveFile([]byte(file.FileID), file.FileName, "./uploads")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to upload file. Please try again.")
				msg.ReplyMarkup = getExitKeyboard()
				bot.Send(msg)
				return
			}
			ad.Files = filePath
			h.userStates[telegramID] = "creating_ad_finish"
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please upload a valid file or press 'Skip'.")
			msg.ReplyMarkup = getExitKeyboard()
			bot.Send(msg)
			return
		}

		h.tempAdData[telegramID] = ad
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Do you want to submit this ad? Use the buttons below:")
		msg.ReplyMarkup = getAdCreationButtons("creating_ad_finish")
		bot.Send(msg)

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
			ad.Approved = false
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
				tgbotapi.NewKeyboardButton("✅ Skip"),
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

func (h *Handler) HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID

	switch data {
	case "add_balance":
		msg := tgbotapi.NewMessage(chatID, "Please enter the amount to add to your balance:")
		bot.Send(msg)
		h.userStates[callbackQuery.From.ID] = "adding_balance"

	case "request_payout":
		msg := tgbotapi.NewMessage(chatID, "Please enter the amount to request for payout:")
		bot.Send(msg)
		h.userStates[callbackQuery.From.ID] = "requesting_payout"

	case "change_name":
		msg := tgbotapi.NewMessage(chatID, "Please enter your new name:")
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
				"ID: %d\nTitle: %s\nPrice: %.2f\nStock: %d\n\n",
				ad.ID, ad.Title, ad.Price, ad.Stock,
			)
		}

		msg := tgbotapi.NewMessage(chatID, adsMessage)
		msg.ParseMode = "Markdown"
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
				"\n*Title:* %d\n*Price:* %s\n*Description:* %s\n\n",
				ad.Title, ad.Price, ad.Description,
			)
		}

		msg := tgbotapi.NewMessage(chatID, ordersMessage)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	case "change_photo":
		msg := tgbotapi.NewMessage(chatID, "Please upload your new profile picture:")
		bot.Send(msg)
		h.userStates[callbackQuery.From.ID] = "changing_photo"
	default:
		msg := tgbotapi.NewMessage(chatID, "Unknown action.")
		bot.Send(msg)
	}
}

func getExitKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Exit"),
		),
	)
}
