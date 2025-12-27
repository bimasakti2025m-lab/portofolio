package controller

import (
	"enigmacamp.com/toko-enigma/model"
	"enigmacamp.com/toko-enigma/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC usecase.AuthenticateUsecase
	rg     *gin.RouterGroup
}

func (a *AuthController) Route() {
	a.rg.POST("/login", a.LoginHandler)
	a.rg.POST("/register", a.RegisterHandler)
}

func (a *AuthController) RegisterHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"err : Binding JSON because": err.Error()})
		return
	}
	user, err := a.authUC.Register(payload)
	if err != nil {
		c.JSON(400, gin.H{"err : Register because": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (a *AuthController) LoginHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"err : Binding JSON because": err.Error()})
		return
	}
	token, err := a.authUC.Login(payload.Username, payload.Password)
	if err != nil {
		c.JSON(400, gin.H{"err : Login because": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func NewAuthController(authUC usecase.AuthenticateUsecase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUC: authUC, rg: rg}
}