package controller

import (
	"net/http"

	"log" // Import package log

	"enigmacamp.com/url-shortener/middleware"
	"enigmacamp.com/url-shortener/model"
	"enigmacamp.com/url-shortener/usecase"
	modelutils "enigmacamp.com/url-shortener/utils/model_utils"
	"github.com/gin-gonic/gin"
)

type UrlController struct {
	useCase usecase.UrlUsecase
	engine  *gin.Engine
	authMid middleware.AuthMiddleware
}

func (u *UrlController) Route() {
	// Endpoint untuk membuat short URL (terproteksi)
	apiGroup := u.engine.Group("/api/v1")
	apiGroup.POST("/urls", u.authMid.RequireToken("admin", "user"), u.createShortUrl)

	// Endpoint untuk redirect (publik)
	// Diletakkan di root router agar URL-nya pendek, misal: http://localhost:8080/xY2zAbC
	u.engine.GET("/:shortCode", u.redirectToLongUrl)
}

func (u *UrlController) createShortUrl(c *gin.Context) {
	var payload model.Url
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user ID dari token JWT yang sudah divalidasi oleh middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: claims not found"})
		return
	}
	payload.UserId = claims.(*modelutils.JwtPayloadClaims).UserId

	createdUrl, err := u.useCase.CreateShortUrl(payload)
	if err != nil {
		log.Printf("Error creating short URL: %v", err) // Log error yang sebenarnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	c.JSON(http.StatusCreated, createdUrl)
}

func (u *UrlController) redirectToLongUrl(c *gin.Context) {
	shortCode := c.Param("shortCode")
	longUrl, err := u.useCase.GetLongUrl(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusMovedPermanently, longUrl)
}

func NewUrlController(useCase usecase.UrlUsecase, engine *gin.Engine, authMid middleware.AuthMiddleware) *UrlController {
	return &UrlController{useCase: useCase, engine: engine, authMid: authMid}
}
