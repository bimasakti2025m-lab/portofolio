package middleware

import (
	"basic-JWT/utils/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	RequireToken func(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JWTservice
}

func (a *authMiddleware) requireToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(401, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
		// trim bearer
		tokenString = tokenString[7:]

		claims, err := a.jwtService.VerifyToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{
			"message": "forbidden",
		})
		c.Abort()
	}
}

func NewAuthMiddleware(jwtService service.JWTservice) *AuthMiddleware {
	am := &authMiddleware{
		jwtService: jwtService,
	}
	return &AuthMiddleware{RequireToken: am.requireToken}
}
