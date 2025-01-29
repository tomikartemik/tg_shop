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
		fmt.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(paymentCallback)

	err := h.services.ChangeStatus(paymentCallback.OrderID, paymentCallback.Status)
	if err != nil {
		fmt.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}
