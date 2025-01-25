package service

import (
	"bytes"
	"encoding/json"
	"fmt"
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
}

func NewCryptoCloudService(repoUser repository.User, repoInvoice repository.Invoice) *CryptoCloudService {
	apiToken = os.Getenv("API_TOKEN")
	apiToken = os.Getenv("SHOP_ID")
	return &CryptoCloudService{
		repoUser:    repoUser,
		repoInvoice: repoInvoice,
	}
}

func (s *CryptoCloudService) CreateInvoice(amount float64, telegramID int) (string, error) {
	orderID, err := s.repoInvoice.CreateInvoice(telegramID, amount)

	if err != nil {
		return "", err
	}

	url := "https://api.cryptocloud.plus/v2/invoice/create"

	// Формируем тело запроса
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

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create invoice: status code %d", resp.StatusCode)
	}

	// Декодируем ответ
	var response InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Проверка, есть ли данные в ответе
	if response.Status != "success" {
		return "", fmt.Errorf("failed to create invoice: status is %s", response.Status)
	}

	return response.Result.Link, nil
}
