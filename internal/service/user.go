package service

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"tg_shop/internal/model"
	"tg_shop/internal/repository"
	"tg_shop/utils"
	"time"
)

type UserService struct {
	repo         repository.User
	repoAd       repository.Ad
	repoCategory repository.Category
	repoEarning  repository.Earning
	bot          *tgbotapi.BotAPI
}

func NewUserService(repo repository.User, repoAd repository.Ad, repoCategory repository.Category, repoEarning repository.Earning, bot *tgbotapi.BotAPI) *UserService {
	return &UserService{
		repo:         repo,
		repoAd:       repoAd,
		repoCategory: repoCategory,
		repoEarning:  repoEarning,
		bot:          bot,
	}
}

func (s *UserService) CreateUser(id int, user model.User) (model.User, error) {
	newUser, err := s.repo.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil
}

func (s *UserService) GetUserInfoById(id int) (model.UserInfo, error) {
	userInfo, err := s.convertUserToUserInfo(id)
	if err != nil {
		return model.UserInfo{}, err
	}

	return userInfo, nil
}

func (s *UserService) GetUserById(id int) (model.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *UserService) CreateOrUpdateUser(user model.User) (model.User, error) {
	existingUser, err := s.repo.GetUserById(user.TelegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			createdUser, createErr := s.repo.CreateUser(user)
			if createErr != nil {
				return model.User{}, createErr
			}
			return createdUser, nil
		}
		return model.User{}, err
	}

	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.PhotoURL != "" {
		existingUser.PhotoURL = user.PhotoURL
	}

	updatedUser, updateErr := s.repo.UpdateUser(existingUser)
	if updateErr != nil {
		return model.User{}, updateErr
	}

	return updatedUser, nil
}

func (s *UserService) UpdateUserName(telegramID int, username string) error {
	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return err
	}

	user.Username = username
	_, updateErr := s.repo.UpdateUser(user)
	return updateErr
}

func (s *UserService) GetUserAsSellerByID(telegramIDStr string) (model.UserAsSeller, error) {
	userAsSeller := model.UserAsSeller{}

	telegramID, err := strconv.Atoi(telegramIDStr)

	if err != nil {
		return model.UserAsSeller{}, err
	}

	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return userAsSeller, err
	}

	var activeAds []model.Ad

	for _, ad := range user.Ads {
		if ad.Status == "Enabled" {
			activeAds = append(activeAds, ad)
		}
	}

	adsShortInfo, err := s.convertAdsToAdsShortInfo(activeAds)

	if err != nil {
		return userAsSeller, err
	}

	userAsSeller = model.UserAsSeller{
		TelegramID:   user.TelegramID,
		Username:     user.Username,
		PhotoURL:     user.PhotoURL,
		Ads:          adsShortInfo,
		Rating:       user.Rating,
		ReviewNumber: user.ReviewNumber,
	}

	return userAsSeller, nil
}

func (s *UserService) IsAdmin(userID int) (bool, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

func (s *UserService) BroadcastMessage(message string) error {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return fmt.Errorf("failed to get all users: %w", err)
	}

	for _, user := range users {
		if err := s.SendMessageToUser(user.TelegramID, message); err != nil {
			log.Printf("Failed to send message to user %d: %v", user.TelegramID, err)
		}
	}

	return nil
}

func (s *UserService) BroadcastAboutDelete(sellerID int, message string) error {
	user, err := s.repo.GetUserById(sellerID)
	if err != nil {
		return fmt.Errorf("failed to get all user: %w", err)
	}

	if err2 := s.SendMessageToUser(user.TelegramID, message); err2 != nil {
		log.Printf("Failed to send message to user %d: %v", user.TelegramID, err2)
	}

	return nil
}

func (s *UserService) SendMessageToUser(telegramID int, message string) error {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("telegram bot token not provided")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	msg := tgbotapi.NewMessage(int64(telegramID), message)
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to user %d: %w", telegramID, err)
	}

	return nil
}

func (s *UserService) BlockUser(userID int) error {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	user.Banned = true
	_, updateErr := s.repo.UpdateUser(user)
	if updateErr != nil {
		return fmt.Errorf("failed to block user: %w", updateErr)
	}

	return nil
}

func (s *UserService) GrantPremium(userID int) error {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	user.IsPremium = true
	user.ExpirePremium = time.Now().AddDate(0, 1, 0)
	_, updateErr := s.repo.UpdateUser(user)
	if updateErr != nil {
		return fmt.Errorf("failed to grant premium: %w", updateErr)
	}

	return s.repoAd.EnableAllDisabledAds(userID)
}

func (s *UserService) ChangeBalance(userID int, newBalance float64) error {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	user.Balance = utils.RoundToTwoDecimalPlaces(newBalance)
	_, updateErr := s.repo.UpdateUser(user)
	if updateErr != nil {
		return fmt.Errorf("failed to change balance: %w", updateErr)
	}

	return nil
}

func (s *UserService) GetUserByUsername(username string) (model.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *UserService) SearchUsers(query string) ([]model.UserInfo, error) {
	var usersInfo []model.UserInfo
	users, err := s.repo.SearchUsers(query)
	if err != nil {
		return usersInfo, err
	}

	for _, user := range users {
		userInfo, err := s.convertUserToUserInfo(user.TelegramID)
		if err != nil {
			continue
		}
		usersInfo = append(usersInfo, userInfo)
	}

	return usersInfo, nil
}

func (s *UserService) ChangeRatingAdm(userID int, newRating float64) error {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	user.Rating = newRating
	_, updateErr := s.repo.UpdateUser(user)
	if updateErr != nil {
		return fmt.Errorf("failed to change rating: %w", updateErr)
	}

	return nil
}

func (s *UserService) Purchase(request model.PurchaseRequest) error {
	userID, adID := request.UserID, request.AdID

	buyer, err := s.repo.GetUserById(userID)
	if err != nil {
		return err
	}

	ad, err := s.repoAd.GetAdById(adID)
	if err != nil {
		return err
	}

	seller, err := s.repo.GetUserById(ad.SellerID)
	if err != nil {
		return err
	}

	if buyer.TelegramID == seller.TelegramID {
		return fmt.Errorf("You really want to buy your own merchandise?")
	}

	if buyer.Balance < ad.Price {
		return errors.New("The financial situation doesn't match!")
	}

	if ad.Stock <= 0 {
		return errors.New("Out of stock!")
	}

	if err = s.repo.AddPurchase(userID, adID); err != nil {
		return err
	}

	var priceForSeller float64

	if seller.IsPremium == true {
		priceForSeller = ad.Price * 0.95
	} else {
		salesCount, err := s.repoEarning.CountEarningsById(seller.TelegramID)

		if err != nil {
			return err
		}

		switch {
		case salesCount >= 40:
			priceForSeller = ad.Price * 0.94
		case salesCount >= 10:
			priceForSeller = ad.Price * 0.93
		default:
			priceForSeller = ad.Price * 0.92
		}
	}

	priceForSeller = (priceForSeller * 100) / 100

	sellerNewHoldBalance := seller.HoldBalance + priceForSeller
	buyerNewBalance := buyer.Balance - ad.Price

	priceForSeller = utils.RoundToTwoDecimalPlaces(priceForSeller)
	sellerNewHoldBalance = utils.RoundToTwoDecimalPlaces(sellerNewHoldBalance)
	buyerNewBalance = utils.RoundToTwoDecimalPlaces(buyerNewBalance)

	if err = s.repo.ChangeHoldBalance(seller.TelegramID, sellerNewHoldBalance); err != nil {
		return err
	}

	newEarning := model.Earning{
		SellerID: seller.TelegramID,
		BuyerID:  buyer.TelegramID,
		Amount:   priceForSeller,
		Status:   "Processing",
	}

	err = s.repoEarning.CreateEarning(newEarning)

	if err != nil {
		return err
	}

	if err = s.repo.ChangeBalance(userID, buyerNewBalance); err != nil {
		return err
	}

	if err = s.repoAd.ChangeStock(adID, ad.Stock-1); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(int64(userID), fmt.Sprintf("Your '%s' purchase has been successfully completed", ad.Title))
	s.bot.Send(msg)

	sellerMsg := fmt.Sprintf("A copy of your ad \"%s\" has just been purchased💰\n\nRemaining number of copies: %d\n\n❗️Money will be available for withdrawal after the inspection time has expired", ad.Title, ad.Stock-1)
	if err := s.SendMessageToUser(seller.TelegramID, sellerMsg); err != nil {
		log.Printf("Failed to send message to seller %d: %v", seller.TelegramID, err)
	}

	if ad.Files != "" {

		absPath, err := filepath.Abs(filepath.Join("..", "cmd", ad.Files))
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		file, err := os.Open(absPath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Получаем расширение файла
		fileExt := filepath.Ext(ad.Files)

		// Формируем имя файла
		fileName := fmt.Sprintf("%s%s", ad.Title, fileExt)

		// Создаем объект для отправки файла
		document := tgbotapi.NewDocument(int64(userID), tgbotapi.FileReader{
			Name:   fileName,
			Reader: file,
		})

		// Отправляем файл через бота
		if _, err := s.bot.Send(document); err != nil {
			return fmt.Errorf("failed to send document: %w", err)
		}

	}

	return nil
}

func (s *UserService) ChangeRating(sellerID int, review int) error {
	user, err := s.repo.GetUserById(sellerID)

	if err != nil {
		return err
	}

	newUserRating := (user.Rating*float64(user.ReviewNumber) + float64(review)) / float64(user.ReviewNumber+1)
	newUserRating = math.Round(newUserRating*100) / 100

	user.Rating = newUserRating
	user.ReviewNumber = user.ReviewNumber + 1

	s.repo.UpdateUser(user)

	return nil
}

func (s *UserService) convertAdsToAdsShortInfo(ads []model.Ad) ([]model.AdShortInfo, error) {
	adsShortInfo := []model.AdShortInfo{}

	for _, ad := range ads {
		seller, err := s.repo.GetUserById(ad.SellerID)

		if err != nil {
			return adsShortInfo, err
		}

		category, err := s.repoCategory.GetCategoryById(ad.CategoryID)

		if err != nil {
			return adsShortInfo, err
		}

		adsShortInfo = append(adsShortInfo, model.AdShortInfo{
			ID:                 ad.ID,
			Title:              ad.Title,
			Price:              ad.Price,
			PhotoURL:           ad.PhotoURL,
			CategoryID:         category.ID,
			CategoryName:       category.Name,
			SellerID:           seller.TelegramID,
			SellerName:         seller.Username,
			SellerRating:       seller.Rating,
			SellerReviewNumber: seller.ReviewNumber,
			Stock:              ad.Stock,
		})
	}

	return adsShortInfo, nil
}

func (s *UserService) convertUserToUserInfo(telegramID int) (model.UserInfo, error) {
	var userInfo model.UserInfo
	user, err := s.repo.GetUserById(telegramID)
	if err != nil {
		return userInfo, err
	}

	purchased, err := s.convertAdsToAdsShortInfo(user.Purchased)
	if err != nil {
		return userInfo, err
	}

	var activeAds []model.Ad

	for _, ad := range user.Ads {
		if ad.Status == "Enabled" {
			activeAds = append(activeAds, ad)
		}
	}

	ads, err := s.convertAdsToAdsShortInfo(activeAds)
	if err != nil {
		return userInfo, err
	}

	userInfo = model.UserInfo{
		TelegramID:   user.TelegramID,
		Username:     user.Username,
		PhotoURL:     user.PhotoURL,
		Balance:      utils.RoundToTwoDecimalPlaces(user.Balance),
		Ads:          ads,
		Purchased:    purchased,
		Rating:       user.Rating,
		ReviewNumber: user.ReviewNumber,
	}

	return userInfo, nil
}
