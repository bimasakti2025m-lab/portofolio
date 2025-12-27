package middleware

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"enigmacamp.com/livecode-catatan-keuangan/shared/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authHeader AuthHeader
		if err := ctx.ShouldBindHeader(&authHeader); err != nil {
			log.Printf("RequireToken: Error binding header: %v \n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println("AuthHeader: ", authHeader)

		tokenHeader := strings.TrimPrefix(authHeader.AuthorizationHeader, "Bearer ")
		if tokenHeader == "" {
			log.Println("RequireToken: Missing token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println("TokenHeader: ", tokenHeader)

		// claims, err := a.jwtService.ParseToken(tokenHeader)
		// if err != nil {
		// 	log.Printf("RequireToken: Error parsing token: %v \n", err)
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		// ctx.Set("user", claims["userId"])

		// role, ok := claims["role"]
		// if !ok {
		// 	log.Println("RequireToken: Missing role in token")
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		// TODO: Encode token base64 to String

		tokenBytes, err := base64.StdEncoding.DecodeString(tokenHeader)
		if err != nil {
			log.Printf("RequireToken: Error decoding token: %v \n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenArr := strings.Split(string(tokenBytes), ":")

		fmt.Println("Split token: ", tokenArr)

		// TODO: Get role by Header
		// [user, ed0f5289-e913-43b3-a0db-78b3ea1f39c7]
		role := tokenArr[0]

		if !isValidRole(role, roles) {
			log.Println("RequireToken: Invalid role")
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// TODO: Set userID to context
		ctx.Set("user", tokenArr[1])

		ctx.Next()
	}
}

func isValidRole(userRole string, validRoles []string) bool {
	for _, role := range validRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
