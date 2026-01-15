package controller

import (
	"net/http"

	"enigmacamp.com/golang-jwt/middleware"
	"enigmacamp.com/golang-jwt/model"
	"enigmacamp.com/golang-jwt/usecase"
	"github.com/gin-gonic/gin"
)

// Mendeklarasikan struct
type TodosController struct {
	todosUseCase usecase.TodosUseCase
	rg           *gin.RouterGroup
	authMid      middleware.AuthMiddleware
}

// Mendeklarasikan endpoint
func (t *TodosController) Route() {
	t.rg.GET("/todos", t.authMid.RequireToken("admin", "user"), t.List)
	t.rg.GET("/todos/:id", t.authMid.RequireToken("admin", "user"), t.Get)
	t.rg.POST("/todos", t.authMid.RequireToken("admin"), t.Create)
	t.rg.PUT("/todos/:id", t.authMid.RequireToken("admin"), t.Update)
	t.rg.DELETE("/todos/:id", t.authMid.RequireToken("admin"), t.Delete)
}

// Implementasi dari interface
func (t *TodosController) Create(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	todo, err := t.todosUseCase.Create(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func (t *TodosController) List(c *gin.Context) {
	todos, err := t.todosUseCase.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (t *TodosController) Get(c *gin.Context) {
	id := c.Param("id")
	todo, err := t.todosUseCase.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (t *TodosController) Update(c *gin.Context) {
	id := c.Param("id")
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	todo, err := t.todosUseCase.Update(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (t *TodosController) Delete(c *gin.Context) {
	id := c.Param("id")
	err := t.todosUseCase.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// Mendeklarasikan konstruktor
func NewTodosController(todosUseCase usecase.TodosUseCase, rg *gin.RouterGroup, authMid middleware.AuthMiddleware) *TodosController {
	return &TodosController{todosUseCase: todosUseCase, rg: rg, authMid: authMid}
}
