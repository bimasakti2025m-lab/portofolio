package controller

import (
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productUc usecase.ProductUsecase
	rg        *gin.RouterGroup
	am        *middleware.AuthMiddleware
}

func (pc *ProductController) Route() {
	pc.rg.POST("/products", pc.am.RequireToken("admin"), pc.createProductHandler)
	pc.rg.GET("/products", pc.am.RequireToken("admin"), pc.getAllProductHandler)
	pc.rg.GET("/products/:id", pc.am.RequireToken("admin"), pc.getProductByIdHandler)
	pc.rg.PUT("/products/:id", pc.am.RequireToken("admin"), pc.updateProductHandler)
	pc.rg.DELETE("/products/:id", pc.am.RequireToken("admin"), pc.deleteProductHandler)
}

func (pc *ProductController) createProductHandler(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	createdProduct, err := pc.productUc.CreateProduct(product)
	if err != nil {
		c.JSON(500, gin.H{"error create product": err.Error()})
		return
	}
	c.JSON(201, createdProduct)
}

func (pc *ProductController) getAllProductHandler(c *gin.Context) {
	products, err := pc.productUc.GetAllProducts()
	if err != nil {
		c.JSON(500, gin.H{"error get product": err.Error()})
		return
	}
	c.JSON(200, products)
}

func (pc *ProductController) getProductByIdHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error convert": err.Error()})
		return
	}

	product, err := pc.productUc.GetProductByID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error get product": err.Error()})
		return
	}
	c.JSON(200, product)
}

func (pc *ProductController) updateProductHandler(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error binding": err.Error()})
		return
	}
	updatedProduct, err := pc.productUc.UpdateProduct(&product)
	if err != nil {
		c.JSON(500, gin.H{"error update product": err.Error()})
		return
	}
	c.JSON(200, updatedProduct)
}

func (pc *ProductController) deleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error convert": err.Error()})
		return
	}
	err = pc.productUc.DeleteProduct(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error delete product": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}

func NewProductController(rg *gin.RouterGroup, productUc usecase.ProductUsecase, am *middleware.AuthMiddleware) *ProductController {
	return &ProductController{productUc: productUc, rg: rg, am: am}
}
