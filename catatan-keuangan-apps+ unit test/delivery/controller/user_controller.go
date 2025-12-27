package controller

import (
	"net/http"

	"enigmacamp.com/livecode-catatan-keuangan/delivery/middleware"
	"enigmacamp.com/livecode-catatan-keuangan/shared/common"
	"enigmacamp.com/livecode-catatan-keuangan/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUc  usecase.UserUseCase
	rg      *gin.RouterGroup
	authMid middleware.AuthMiddleware
}

func (u *UserController) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := u.userUc.FindUserByID(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}
	common.SendSingleResponse(ctx, user, "Ok")
}

func (u *UserController) Route() {
	u.rg.GET("/users/:id", u.authMid.RequireToken("user", "admin"), u.getHandler)
}

func NewUserController(userUc usecase.UserUseCase, rg *gin.RouterGroup, authMid middleware.AuthMiddleware) *UserController {
	return &UserController{userUc: userUc, rg: rg, authMid: authMid}
}
