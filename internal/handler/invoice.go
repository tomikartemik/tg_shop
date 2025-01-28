package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tg_shop/internal/model"
)

func (h *Handler) paymentCallback(c *gin.Context) {
	var paymentCallback model.PaymentCallback
	if err := c.ShouldBindJSON(&paymentCallback); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(paymentCallback)

	err := h.services.ChangeStatus(paymentCallback.OrderID, paymentCallback.Status)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Покупка успешно совершена"})
}
