package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tg_shop/internal/model"
	"tg_shop/internal/service"
)

type Handler struct {
	services   *service.Service
	userStates map[int64]string
	tempAdData map[int64]model.Ad
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services:   services,
		userStates: make(map[int64]string),
		tempAdData: make(map[int64]model.Ad),
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

	//USER
	////////////////////////////////////////////////////////////
	user := router.Group("/user")
	{
		user.GET("", h.GetUserByID)
		user.GET("/seller", h.GetUserAsSellerByID)
	}
	////////////////////////////////////////////////////////////

	//AD
	////////////////////////////////////////////////////////////
	ad := router.Group("/ad")
	{
		ad.GET("/list", h.GetAdList)
		ad.GET("", h.GetAdById)
	}
	////////////////////////////////////////////////////////////

	//CATEGORY
	////////////////////////////////////////////////////////////
	category := router.Group("/category")
	{
		category.GET("", h.GetCategoryList)
	}
	////////////////////////////////////////////////////////////

	return router
}
