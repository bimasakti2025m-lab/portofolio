package controller

import (

	"net/http"
	"strconv"

	"enigmacamp.com/toko-enigma/middleware"
	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/usecase"
	"github.com/gin-gonic/gin"
)

// BookController adalah sebuah struct yang menghandle semua operasi terkait buku
type CartController struct {
	useCase        usecase.CartUseCase       // use case untuk operasi buku
	rg             *gin.RouterGroup          // router group untuk menghandle request
	authMiddleware middleware.AuthMiddleware // middleware untuk autentikasi
}

// Route adalah sebuah method yang mengatur routing untuk operasi buku
func (c *CartController) Route() {
	c.rg.POST("/carts", c.authMiddleware.RequireToken("admin"), c.createCart)
	c.rg.GET("/carts", c.authMiddleware.RequireToken("admin", "user"), c.getAllCarts)
	c.rg.GET("/carts/:id", c.authMiddleware.RequireToken("admin", "user"), c.getCartById)
	c.rg.PUT("/carts", c.authMiddleware.RequireToken("admin"), c.updateCart)
	c.rg.DELETE("/carts/:id", c.authMiddleware.RequireToken("admin"), c.deleteCart)
}

func (b *CartController) createCart(c *gin.Context) {
	var payload model.Cart
	if err := c.ShouldBindJSON(&payload); err != nil {
		// jika terjadi error saat binding data, maka return error
		c.JSON(http.StatusBadRequest, gin.H{"err": "Failed to bind JSON" + "because " + err.Error()})
		return
	}

	cart, err := b.useCase.CreateNewCart(payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create cart"+ "because " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cart)
}

func (b *CartController) getAllCarts(c *gin.Context) {
	carts, err := b.useCase.GetAllCart()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data cart"+ "because " + err.Error()})
		return
	}

	if len(carts) > 0 {
		c.JSON(http.StatusOK, carts)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List cart empty"})
}

func (b *CartController) getCartById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cart, err := b.useCase.GetCartById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get cart by ID"+ "because " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (b *CartController) updateCart(c *gin.Context) {
	var payload model.Cart
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Failed to bind JSON" + "because " + err.Error()})
		return
	}

	cart, err := b.useCase.UpdateCartById(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to update cart"+ "because " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (b *CartController) deleteCart(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := b.useCase.DeleteCartById(id)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to delete todo"+ "because " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Cart deleted successfully.")
}

func NewCartController(useCase usecase.CartUseCase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *CartController {
	return &CartController{useCase: useCase, rg: rg, authMiddleware: am}
}
