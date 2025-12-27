package controller

import (
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartUc usecase.CartUseCase
	rg     *gin.RouterGroup
	am     *middleware.AuthMiddleware
}

func (cc *CartController) Route() {
	cc.rg.POST("/carts", cc.am.RequireToken("user"), cc.createCartHandler)
	cc.rg.GET("/carts", cc.am.RequireToken("user"), cc.getAllCartsHandler)
	cc.rg.GET("/carts/:id", cc.am.RequireToken("user"), cc.getCartByIdHandler)
	cc.rg.PUT("/carts/:id", cc.am.RequireToken("user"), cc.updateCartHandler)
	cc.rg.DELETE("/carts/:id", cc.am.RequireToken("user"), cc.deleteCartHandler)
}

func (cc *CartController) createCartHandler(c *gin.Context) {
	var cart model.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	createdCart, err := cc.cartUc.CreateCart(cart)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdCart)
}

func (cc *CartController) getAllCartsHandler(c *gin.Context) {
	carts, err := cc.cartUc.GetAllCarts()
	if err != nil {
		c.JSON(500, gin.H{"error get cart": err.Error()})
		return
	}

	c.JSON(200, carts)
}

func (cc *CartController) getCartByIdHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error convert": err.Error()})
		return
	}

	cart, err := cc.cartUc.GetCartByID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error get cart": err.Error()})
		return
	}

	c.JSON(200, cart)
}

func (cc *CartController) updateCartHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error convert": err.Error()})
		return
	}
	var cart model.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	cart.ID = idInt
	updatedCart, err := cc.cartUc.UpdateCart(cart)
	if err != nil {
		c.JSON(500, gin.H{"error update cart": err.Error()})
		return
	}

	c.JSON(200, updatedCart)
}

func (cc *CartController) deleteCartHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error convert": err.Error()})
		return
	}
	err = cc.cartUc.DeleteCart(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error delete cart": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart deleted"})
}

func NewCartController(rg *gin.RouterGroup, cartUc usecase.CartUseCase, am *middleware.AuthMiddleware) *CartController {
	return &CartController{cartUc: cartUc, rg: rg, am: am}
}