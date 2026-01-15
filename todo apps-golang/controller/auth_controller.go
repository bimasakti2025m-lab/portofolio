package controller

import (
	"net/http"

	"enigmacamp.com/golang-jwt/model"
	"enigmacamp.com/golang-jwt/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC usecase.AuthenticateUsecase
	rg     *gin.RouterGroup
}

func (a *AuthController) Route() {
	if a.rg == nil {
		panic("AuthController: RouterGroup is nil. Check your initialization in main/server.")
	}
	a.rg.POST("/login", a.LoginHandler)
	a.rg.POST("/register", a.RegisterHandler)
	a.rg.GET("/health", a.HealthHandler)
}

func (a *AuthController) LoginHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	token, err := a.authUC.Login(payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthController) RegisterHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := a.authUC.Register(payload)
	if err != nil {
		// Handle specific errors like duplicate username if needed, otherwise generic 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (a *AuthController) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func NewAuthController(authUC usecase.AuthenticateUsecase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUC: authUC, rg: rg}
}
