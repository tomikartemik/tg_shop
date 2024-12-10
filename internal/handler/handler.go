package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tg_shop/internal/service"
)

type Handler struct {
	services   *service.Service
	userStates map[int64]string // Telegram ID -> состояние или язык
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services:   services,
		userStates: make(map[int64]string),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return router
}
