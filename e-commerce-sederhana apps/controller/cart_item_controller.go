package controller

import (
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartItemController struct {
	cartItemUc usecase.CartItemUseCase
	rg         *gin.RouterGroup
	am         *middleware.AuthMiddleware
}

func (cc *CartItemController) Route() {
	cc.rg.POST("/cart-items", cc.am.RequireToken("user"), cc.createCartItemHandler)
	cc.rg.GET("/cart-items", cc.am.RequireToken("user"), cc.getAllCartItemsHandler)
	cc.rg.GET("/cart-items/:id", cc.am.RequireToken("user"), cc.getCartItemByIdHandler)
	cc.rg.PUT("/cart-items/:id", cc.am.RequireToken("user"), cc.updateCartItemHandler)
	cc.rg.DELETE("/cart-items/:id", cc.am.RequireToken("user"), cc.deleteCartItemHandler)
}

func (cc *CartItemController) createCartItemHandler(c *gin.Context) {
	var cartItem model.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	createdCartItem, err := cc.cartItemUc.CreateCartItem(cartItem)
	if err != nil {
		c.JSON(500, gin.H{"error creating": err.Error()})
		return
	}

	c.JSON(201, createdCartItem)
}

func (cc *CartItemController) getAllCartItemsHandler(c *gin.Context) {
	cartItems, err := cc.cartItemUc.GetAllCartItems()
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}

	c.JSON(200, cartItems)
}

func (cc *CartItemController) getCartItemByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}

	cartItem, err := cc.cartItemUc.GetCartItemByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error getting": err.Error()})
		return
	}

	c.JSON(200, cartItem)
}

func (cc *CartItemController) updateCartItemHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	var cartItem model.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	cartItem.ID = id
	updatedCartItem, err := cc.cartItemUc.UpdateCartItem(cartItem)
	if err != nil {
		c.JSON(500, gin.H{"error updating": err.Error()})
		return
	}

	c.JSON(200, updatedCartItem)
}

func (cc *CartItemController) deleteCartItemHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error converting": err.Error()})
		return
	}
	err = cc.cartItemUc.DeleteCartItem(id)
	if err != nil {
		c.JSON(500, gin.H{"error deleting": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart item deleted"})
}

func NewCartItemController(cartItemUc usecase.CartItemUseCase, rg *gin.RouterGroup, am *middleware.AuthMiddleware) *CartItemController {
	return &CartItemController{cartItemUc: cartItemUc, rg: rg, am: am}
}