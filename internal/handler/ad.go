package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAdList(c *gin.Context) {
	categoryIDStr := c.Query("category_id")

	ads, err := h.services.GetAdList(categoryIDStr)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, ads)

}

func (h *Handler) GetAdById(c *gin.Context) {
	idStr := c.Query("id")

	ad, err := h.services.GetAdByID(idStr)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, ad)
}
