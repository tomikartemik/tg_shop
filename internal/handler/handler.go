package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tg_shop/internal/service"
)

type Handler struct {
	services         *service.Service
	pendingUsernames map[int64]bool // Telegram ID -> ждём имя
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services,
		pendingUsernames: make(map[int64]bool)}
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
