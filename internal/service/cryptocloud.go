package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strconv"
	"tg_shop/internal/repository"
)

var apiToken string
var shopID string

type InvoiceResult struct {
	UUID       string  `json:"uuid"`
	Created    string  `json:"created"`
	ExpiryDate string  `json:"expiry_date"`
	Amount     float64 `json:"amount"`
	AmountUSD  float64 `json:"amount_usd"`
	Link       string  `json:"link"`
	Status     string  `json:"status"`
}

// InvoiceResponse описывает полный ответ от API CryptoCloud
type InvoiceResponse struct {
	Status string        `json:"status"`
	Result InvoiceResult `json:"result"`
}

type CryptoCloudService struct {
	repoUser    repository.User
	repoInvoice repository.Invoice
	bot         *tgbotapi.BotAPI
}

func NewCryptoCloudService(repoUser repository.User, repoInvoice repository.Invoice, bot *tgbotapi.BotAPI) *CryptoCloudService {
	apiToken = os.Getenv("API_TOKEN")
	shopID = os.Getenv("SHOP_ID")
	return &CryptoCloudService{
		repoUser:    repoUser,
		repoInvoice: repoInvoice,
		bot:         bot,
	}
}

func (s *CryptoCloudService) CreateInvoice(amount float64, telegramID int) (string, error) {
	orderID, err := s.repoInvoice.CreateInvoice(telegramID, amount)

	if err != nil {
		return "", err
	}

	url := "https://api.cryptocloud.plus/v2/invoice/create"

	payload := map[string]interface{}{
		"shop_id":  shopID,
		"amount":   fmt.Sprintf("%.2f", amount),
		"order_id": strconv.Itoa(orderID),
	}

	// Преобразуем тело запроса в JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Token "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create invoice: status code %d", resp.StatusCode)
	}

	var response InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Status != "success" {
		return "", fmt.Errorf("failed to create invoice: status is %s", response.Status)
	}

	return response.Result.Link, nil
}

func (s *CryptoCloudService) ChangeStatus(idStr string, status string) error {
	id, err := strconv.Atoi(idStr)

	if status != "success" {
		return s.repoInvoice.ChangeStatus(id, status)
	}

	invoice, err := s.repoInvoice.GetInvoiceByID(id)
	if err != nil {
		return err
	}

	user, err := s.repoUser.GetUserById(invoice.TelegramID)
	if err != nil {
		return err
	}

	newBalance := user.Balance + invoice.Amount

	err = s.repoUser.ChangeBalance(user.TelegramID, newBalance)
	if err != nil {
		return err
	}

	err = s.repoInvoice.ChangeStatus(id, status)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("🎉 Your balance has been successfully topped up by %.2f$. New balance: %.2f$", invoice.Amount, newBalance)
	msg := tgbotapi.NewMessage(int64(user.TelegramID), message)
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("Failed to send notification to user %d: %v", user.TelegramID, err)
	}

	return nil
}
