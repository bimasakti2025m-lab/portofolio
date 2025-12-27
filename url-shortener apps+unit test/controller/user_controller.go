package controller

import (
	"log"
	"net/http"
	"strconv"

	"enigmacamp.com/url-shortener/middleware"
	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	useCase usecase.UserUseCase
	rg      *gin.RouterGroup
	authMid middleware.AuthMiddleware
}

func (u *UserController) Route() {
	u.rg.POST("/users/register",u.authMid.RequireToken("admin"), u.createUser)
	u.rg.GET("/users",u.authMid.RequireToken("admin","user"), u.getAllUser)
	u.rg.GET("/users/:id",u.authMid.RequireToken("admin", "user"), u.getUserById)
}

func (u *UserController) createUser(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := u.useCase.RegisterNewUser(payload)
	if err != nil {
		log.Println("Error creating user:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (u *UserController) getAllUser(c *gin.Context) {
	users, err := u.useCase.FindAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data users"})
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusOK, users)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List user empty"})
}

func (u *UserController) getUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := u.useCase.FindUserById(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get user by ID"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func NewUserController(useCase usecase.UserUseCase, rg *gin.RouterGroup, authMid middleware.AuthMiddleware) *UserController {
	return &UserController{useCase: useCase, rg: rg, authMid: authMid}
}
