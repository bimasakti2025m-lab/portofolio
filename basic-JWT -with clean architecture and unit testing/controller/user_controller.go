package controller

import (
	"basic-JWT/middleware"
	"basic-JWT/model"
	"basic-JWT/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUc         usecase.UserUsecase
	rg             *gin.RouterGroup
	authMiddleware *middleware.AuthMiddleware
}

func (uc *UserController) Route() {
	uc.rg.POST("/users", uc.authMiddleware.RequireToken("admin"), uc.createUserHandler)
	uc.rg.GET("/users", uc.authMiddleware.RequireToken("user", "admin"), uc.getAllUsersHandler)
	uc.rg.GET("/users/:username", uc.authMiddleware.RequireToken("user", "admin"), uc.getUserByUsernameHandler)
}

func (uc *UserController) createUserHandler(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	_, err = uc.userUc.Create(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

func (uc *UserController) getAllUsersHandler(c *gin.Context) {
	users, err := uc.userUc.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to get users",
		})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}

func (uc *UserController) getUserByUsernameHandler(c *gin.Context) {
	username := c.Param("username")
	user, err := uc.userUc.GetUserByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func NewUserController(rg *gin.RouterGroup, userUc usecase.UserUsecase, authMiddleware *middleware.AuthMiddleware) *UserController {
	return &UserController{userUc: userUc, rg: rg, authMiddleware: authMiddleware}
}
