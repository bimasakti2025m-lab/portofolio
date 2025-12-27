package controller

import (
	"log"
	"net/http"
	"strconv"

	"enigmacamp.com/mini-banking/middleware"
	"enigmacamp.com/mini-banking/model"
	"enigmacamp.com/mini-banking/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionUC usecase.TransactionUseCase
	rg            *gin.RouterGroup
	authMid       middleware.AuthMiddleware
}

func (t *TransactionController) Route() {
	t.rg.POST("/transactions", t.authMid.RequireToken("admin", "user"), t.createTransaction)
	t.rg.GET("/transactions", t.authMid.RequireToken("admin"), t.listTransaction)
	t.rg.GET("/transactions/:id", t.authMid.RequireToken("admin", "user"), t.getTransactionById)
	t.rg.GET("/transactions/user/:userId", t.authMid.RequireToken("admin", "user"), t.getTransactionByUserId)
	t.rg.PUT("/transactions", t.authMid.RequireToken("admin"), t.updateTransaction)
	t.rg.DELETE("/transactions/:id", t.authMid.RequireToken("admin"), t.deleteTransaction)
}

func (t *TransactionController) createTransaction(c *gin.Context) {
	var payload model.Transaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	transaction, err := t.transactionUC.CreateTransaction(payload)
	if err != nil {
		log.Println("Error creating transaction:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (t *TransactionController) listTransaction(c *gin.Context) {
	transactions, err := t.transactionUC.ListTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data transactions"})
		return
	}

	if len(transactions) > 0 {
		c.JSON(http.StatusOK, transactions)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List transaction empty"})
}

func (t *TransactionController) getTransactionById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	transaction, err := t.transactionUC.GetTransactionById(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get transaction by ID"})
		return
	}

	c.JSON(http.StatusOK, transaction)

	}

func (t *TransactionController) getTransactionByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	transactions, err := t.transactionUC.GetTransactionByUserId(uint32(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get transaction by user ID"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) updateTransaction(c *gin.Context) {
	var payload model.Transaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	transaction, err := t.transactionUC.UpdateTransaction(payload)
	if err != nil {
		log.Println("Error updating transaction:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (t *TransactionController) deleteTransaction(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := t.transactionUC.DeleteTransaction(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func NewTransactionController(transactionUC usecase.TransactionUseCase, rg *gin.RouterGroup, authMid middleware.AuthMiddleware) *TransactionController {
	return &TransactionController{transactionUC: transactionUC, rg: rg, authMid: authMid}
}