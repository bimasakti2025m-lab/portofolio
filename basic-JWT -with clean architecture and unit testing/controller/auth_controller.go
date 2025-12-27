package controller

import (
	"basic-JWT/model"
	"basic-JWT/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUc usecase.AuthenticationUsecase
	rg     *gin.RouterGroup
}

func (ac *AuthController) Route() {
	ac.rg.POST("/register", ac.registerHandler)
	ac.rg.POST("/login", ac.loginHandler)
}

func (ac *AuthController) registerHandler(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	user, err = ac.authUc.Register(user.Username, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "is already taken") {
			c.JSON(409, gin.H{"error": err.Error()}) // 409 Conflict
			return
		}
		c.JSON(500, gin.H{
			"error": "failed to register user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

func (ac *AuthController) loginHandler(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		return
	}

	token, err := ac.authUc.Login(user.Username, user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to login user",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func NewAuthController(rg *gin.RouterGroup, authUc usecase.AuthenticationUsecase) *AuthController {
	return &AuthController{authUc: authUc, rg: rg}
}
