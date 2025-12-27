// TODO :
// 1. Mendeklarasikan nama package controller
// 2. Mendeklarasikan struct bernama BookController
// 3. Mendeklarasikan funtion route
// 4. Mendeklarasikan detail handler
// 5. Mendeklarasikan function baru bernama newBookController

package controller

import (
	"simple-clean-architecture/model"
	"simple-clean-architecture/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookUsecase usecase.BookUsecase
	rg          *gin.RouterGroup
}

func (b *BookController) Route() {
	b.rg.POST("/books", b.CreateNewBook)
	b.rg.GET("/books", b.GetAllBook)
	b.rg.GET("/books/:id", b.GetBookById)
	b.rg.PUT("/books/:id", b.UpdateBook)
	b.rg.DELETE("/books/:id", b.DeleteBook)
}

func (b *BookController) CreateNewBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	newBook, err := b.bookUsecase.CreateNewBook(book)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, newBook)
}

func (b *BookController) GetAllBook(c *gin.Context) {
	books, err := b.bookUsecase.GetAllBook()

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, books)
}

func (b *BookController) GetBookById(c *gin.Context) {
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	book, err := b.bookUsecase.GetBookById(intID)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, book)
}

func (b *BookController) UpdateBook(c *gin.Context) {
	var book model.Book
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	book.Id = intID

	updatedBook, err := b.bookUsecase.UpdateBook(&book)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, updatedBook)
}

func (b *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	err := b.bookUsecase.DeleteBook(intID)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Book deleted successfully"})
}

func NewBookController(bookUsecase usecase.BookUsecase, rg *gin.RouterGroup) *BookController {
	return &BookController{bookUsecase: bookUsecase, rg: rg}
}
