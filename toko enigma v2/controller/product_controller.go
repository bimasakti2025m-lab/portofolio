package controller

import (
	"net/http"

	"enigmacamp.com/toko-enigma/middleware"
	"enigmacamp.com/toko-enigma/model"

	"strconv"

	"enigmacamp.com/toko-enigma/usecase"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	useCase        usecase.ProductUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (p *ProductController) Route() {
	p.rg.POST("/products", p.authMiddleware.RequireToken("admin"), p.createNewProduct)
	p.rg.GET("/products", p.authMiddleware.RequireToken("admin"), p.getAllProduct)
	p.rg.GET("/products/:id", p.authMiddleware.RequireToken("admin", "user"), p.getProductById)
	p.rg.PUT("/products", p.authMiddleware.RequireToken("admin"), p.updateProductById)
	p.rg.DELETE("/products/:id", p.authMiddleware.RequireToken("admin"), p.deleteProductById)
}

func (p *ProductController) createNewProduct(c *gin.Context) {
	var payload model.Product

	err := c.ShouldBindJSON(&payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	book, err := p.useCase.CreateNewProduct(payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create data product."})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (p *ProductController) getAllProduct(c *gin.Context) {
	books, err := p.useCase.GetAllProduct()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get all data product."})
		return
	}

	if len(books) > 0 {
		c.JSON(http.StatusOK, books)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List product is empty."})
}

func (p *ProductController) getProductById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	book, err := p.useCase.GetProductById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get data product by id."})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (p *ProductController) updateProductById(c *gin.Context) {
	var payload model.Product

	err := c.ShouldBindJSON(&payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	book, err := p.useCase.UpdateProductById(payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update data product by ID with error: " + err.Error() + "."})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (p *ProductController) deleteProductById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := p.useCase.DeleteProductById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to delete product by ID."})
		return
	}

	c.JSON(http.StatusOK, "Product deleted successfully.")
}

func NewProductController(useCase usecase.ProductUseCase, rg *gin.RouterGroup, am middleware.AuthMiddleware) *ProductController {
	return &ProductController{useCase: useCase, rg: rg,authMiddleware: am}
}
