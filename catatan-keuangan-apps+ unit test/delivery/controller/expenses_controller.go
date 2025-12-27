package controller

import (
	"net/http"
	"strconv"

	"enigmacamp.com/livecode-catatan-keuangan/config"
	"enigmacamp.com/livecode-catatan-keuangan/delivery/middleware"
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/shared/common"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"enigmacamp.com/livecode-catatan-keuangan/usecase"
	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	expenseUc usecase.ExpenseUseCase
	rg        *gin.RouterGroup
	authMid   middleware.AuthMiddleware
}

func (e *ExpenseController) createHandler(ctx *gin.Context) {
	var payload entity.Expense
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user := ctx.MustGet("user").(string)
	payload.UserId = user
	rsv, err := e.expenseUc.RegisterNewExpense(payload)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	common.SendCreateResponse(ctx, rsv)
}

func (e *ExpenseController) listHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	user := ctx.MustGet("user").(string)

	rsv, paging, err := e.expenseUc.FindAllExpense(page, size, startDate, endDate, user)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var interfaceSlice = make([]interface{}, len(rsv))
	for i, v := range rsv {
		interfaceSlice[i] = v
	}
	common.SendPagedResponse(ctx, interfaceSlice, paging, "Ok")
}

func (e *ExpenseController) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	rsv, err := e.expenseUc.FindExpenseByID(id)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, "not found ID "+id)
		return
	}
	common.SendSingleResponse(ctx, rsv, "Ok")
}

func (e *ExpenseController) getByTransactionHandler(ctx *gin.Context) {
	transactionType := ctx.Param("type")
	user := ctx.MustGet("user").(string)
	rsv, err := e.expenseUc.FindExpenseByTransactionType(transactionType, user)
	if err != nil {
		common.SendErrorResponse(ctx, http.StatusNotFound, "not found type "+transactionType)
		return
	}
	var interfaceSlice = make([]interface{}, len(rsv))
	for i, v := range rsv {
		interfaceSlice[i] = v
	}
	common.SendPagedResponse(ctx, interfaceSlice, model.Paging{}, "Ok")
}

func (e *ExpenseController) Route() {
	e.rg.POST(config.PostExpense, e.authMid.RequireToken("user"), e.createHandler)
	e.rg.GET(config.GetExpenseList, e.authMid.RequireToken("user"), e.listHandler)
	e.rg.GET(config.GetExpense, e.authMid.RequireToken("user"), e.getHandler)
	e.rg.GET(config.GetExpenseTransaction, e.authMid.RequireToken("user"), e.getByTransactionHandler)
}

func NewExpenseController(expenseUc usecase.ExpenseUseCase, rg *gin.RouterGroup, authMid middleware.AuthMiddleware) *ExpenseController {
	return &ExpenseController{expenseUc: expenseUc, rg: rg, authMid: authMid}
}
