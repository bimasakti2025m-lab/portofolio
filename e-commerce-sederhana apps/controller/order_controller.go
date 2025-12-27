package controller

import (
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderUc usecase.OrderUsecase
	rg      *gin.RouterGroup
	am      *middleware.AuthMiddleware
}

func (oc *OrderController) Route() {
	oc.rg.POST("/orders", oc.am.RequireToken("user"), oc.createOrderHandler)
	oc.rg.GET("/orders", oc.am.RequireToken("user"), oc.getAllOrdersHandler)
	oc.rg.GET("/orders/:id", oc.am.RequireToken("user"), oc.getOrderByIdHandler)
	oc.rg.PUT("/orders/:id", oc.am.RequireToken("user"), oc.updateOrderHandler)
	oc.rg.DELETE("/orders/:id", oc.am.RequireToken("user"), oc.deleteOrderHandler)
}

func (oc *OrderController) createOrderHandler(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}

	// Ambil userID dari context yang di-set oleh middleware
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	order.UserID = userID.(int)
	createdOrder, err := oc.orderUc.CreateOrder(order)
	if err != nil {
		c.JSON(500, gin.H{"error creating": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message":     "Order created successfully, please proceed to payment.",
		"payment_url": createdOrder.TransactionIDMidtrans, // Field ini sekarang berisi URL
	})
}

func (oc *OrderController) getAllOrdersHandler(c *gin.Context) {
	orders, err := oc.orderUc.GetAllOrders()
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}
	c.JSON(200, orders)
}

func (oc *OrderController) getOrderByIdHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	order, err := oc.orderUc.GetOrderById(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}
	c.JSON(200, order)
}

func (oc *OrderController) updateOrderHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	order.ID = idInt
	updatedOrder, err := oc.orderUc.UpdateOrder(order)
	if err != nil {
		c.JSON(500, gin.H{"error updating": err.Error()})
		return
	}
	c.JSON(200, updatedOrder)
}

func (oc *OrderController) deleteOrderHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	err = oc.orderUc.DeleteOrder(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error deleting": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Order deleted successfully"})
}

func NewOrderController(rg *gin.RouterGroup, orderUc usecase.OrderUsecase, am *middleware.AuthMiddleware) *OrderController {
	return &OrderController{orderUc: orderUc, rg: rg, am: am}
}
