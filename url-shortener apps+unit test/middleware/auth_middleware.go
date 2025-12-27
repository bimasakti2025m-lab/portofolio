package middleware

import (
	"net/http"
	"strings"

	"enigmacamp.com/url-shortener/utils/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var header authHeader
		if err := context.ShouldBindHeader(&header); err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		bearerToken := header.Authorization
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "invalid token format"})
			return
		}
		token := strings.Replace(bearerToken, "Bearer ", "", 1)

		claims, err := a.jwtService.VerifyToken(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			return
		}

		context.Set("claims", claims)

		if len(roles) > 0 {
			for _, role := range roles {
				if claims.Role == role {
					context.Next()
					return
				}
			}
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "Forbidden: you don't have the right role"})
			return
		}
		context.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwtService,
	}
}
