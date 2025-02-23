package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tg_shop/internal/model"
)

func (h *Handler) GetUserByID(c *gin.Context) {
	userIDStr := c.Query("tg_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID: "+err.Error())
		return
	}

	user, err := h.services.GetUserInfoById(userID)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserAsSellerByID(c *gin.Context) {
	telegramIDStr := c.Query("tg_id")
	userAsSeller, err := h.services.GetUserAsSellerByID(telegramIDStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, userAsSeller)
}

func (h *Handler) SearchUsers(c *gin.Context) {
	query := c.Query("username")
	users, err := h.services.SearchUsers(query)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Purchase(c *gin.Context) {
	var purchaseRequest model.PurchaseRequest
	if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Purchase(purchaseRequest)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase successfully completed!"})
}
