package controller

import (
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderItemController struct {
	orderItemUc usecase.OrderItemUsecase
	rg         *gin.RouterGroup
	am         *middleware.AuthMiddleware
}

func (oc *OrderItemController) Route() {
	oc.rg.POST("/order-items", oc.am.RequireToken("user"), oc.createOrderItemHandler)
	oc.rg.GET("/order-items", oc.am.RequireToken("user"), oc.getAllOrderItemsHandler)
	oc.rg.GET("/order-items/:id", oc.am.RequireToken("user"), oc.getOrderItemByIdHandler)
	oc.rg.PUT("/order-items/:id", oc.am.RequireToken("user"), oc.updateOrderItemHandler)
	oc.rg.DELETE("/order-items/:id", oc.am.RequireToken("user"), oc.deleteOrderItemHandler)
}

func (oc *OrderItemController) createOrderItemHandler(c *gin.Context) {
	var orderItem model.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	createdOrderItem, err := oc.orderItemUc.CreateOrderItem(orderItem)
	if err != nil {
		c.JSON(500, gin.H{"error creating": err.Error()})
		return
	}
	c.JSON(201, createdOrderItem)
}

func (oc *OrderItemController) getAllOrderItemsHandler(c *gin.Context) {
	orderItems, err := oc.orderItemUc.GetAllOrderItems()
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}
	c.JSON(200, orderItems)
}

func (oc *OrderItemController) getOrderItemByIdHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	orderItem, err := oc.orderItemUc.GetOrderItemByID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}
	c.JSON(200, orderItem)
}

func (oc *OrderItemController) updateOrderItemHandler(c *gin.Context) {
	var orderItem model.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	updatedOrderItem, err := oc.orderItemUc.UpdateOrderItem(orderItem)
	if err != nil {
		c.JSON(500, gin.H{"error updating": err.Error()})
		return
	}
	c.JSON(200, updatedOrderItem)
}

func (oc *OrderItemController) deleteOrderItemHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	err = oc.orderItemUc.DeleteOrderItem(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error deleting": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Order item deleted successfully"})
}

func NewOrderItemController(rg *gin.RouterGroup, orderItemUc usecase.OrderItemUsecase, am *middleware.AuthMiddleware) *OrderItemController {
	return &OrderItemController{orderItemUc: orderItemUc, rg: rg, am: am}
}