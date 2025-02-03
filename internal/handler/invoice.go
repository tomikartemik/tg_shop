package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) paymentCallback(c *gin.Context) {
	// Чтение параметров из тела запроса
	status := c.PostForm("status")
	orderID := c.PostForm("order_id")

	// Вывод параметров для отладки
	fmt.Println("Status:", status)
	fmt.Println("Order ID:", orderID)

	if status == "" || orderID == "" {
		newErrorResponse(c, http.StatusBadRequest, "Отсутствуют необходимые параметры")
		return
	}

	// Изменение статуса заказа
	err := h.services.ChangeStatus(orderID, status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}
