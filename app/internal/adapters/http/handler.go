package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/wb-l0/app/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) InitAPI() *gin.Engine {
	router := gin.Default()

	router.Use(CORS)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := router.Group("/api")
	{
		orders := api.Group("/order")
		{
			orders.GET("", h.getOrders)
			orders.GET("/:order_id", h.getOrderByID)
		}
	}

	return router
}
