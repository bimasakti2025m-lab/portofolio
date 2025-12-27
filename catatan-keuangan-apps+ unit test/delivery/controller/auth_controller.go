package controller

import (
	"fmt"
	"net/http"

	"enigmacamp.com/livecode-catatan-keuangan/entity/dto"
	"enigmacamp.com/livecode-catatan-keuangan/shared/common"
	"enigmacamp.com/livecode-catatan-keuangan/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUc usecase.AuthUseCase
	rg     *gin.RouterGroup
}

func (a *AuthController) loginHandler(ctx *gin.Context) {
	var payload dto.AuthRequestDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	token, err := a.authUc.Login(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendCreateResponse(ctx, token)
}

func (a *AuthController) registerHandler(ctx *gin.Context) {
	var payload dto.AuthRequestDto

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("payload trigger: ", payload)
	user, err := a.authUc.Register(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendCreateResponse(ctx, user)
}

func (a *AuthController) Route() {
	a.rg.POST("/auth/login", a.loginHandler)
	a.rg.POST("/auth/register", a.registerHandler)
}

func NewAuthController(authUc usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUc: authUc, rg: rg}
}
