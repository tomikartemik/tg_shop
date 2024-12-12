package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetCategoryList(c *gin.Context) {
	categories, err := h.services.GetCategoryList()

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, categories)

}
