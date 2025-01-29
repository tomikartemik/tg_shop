package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"tg_shop/internal/model"
)

func (h *Handler) paymentCallback(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Выводим тело запроса в консоль
	fmt.Println("Request Body:", string(body))

	var paymentCallback model.PaymentCallback
	if err := c.ShouldBindJSON(&paymentCallback); err != nil {
		fmt.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(paymentCallback)

	err = h.services.ChangeStatus(paymentCallback.OrderID, paymentCallback.Status)
	if err != nil {
		fmt.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Покупка успешно совершена"})
}
