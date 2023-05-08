package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shamank/wb-l0/app/internal/domain"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) getOrderByID(c *gin.Context) {
	var OrderID string = c.Param("order_id")

	print(OrderID)

	order, err := h.service.Orders.Get(c.Request.Context(), OrderID)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

type getAllOrders struct {
	Orders []domain.Order `json:"orders"`
}

func (h *Handler) getOrders(c *gin.Context) {
	orders, err := h.service.Orders.GetAll(c.Request.Context())
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("orders: %x", orders)

	c.JSON(http.StatusOK, getAllOrders{Orders: orders})
}
