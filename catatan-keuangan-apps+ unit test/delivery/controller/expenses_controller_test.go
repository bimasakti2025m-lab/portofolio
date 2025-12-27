package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enigmacamp.com/livecode-catatan-keuangan/delivery/middleware"
	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ExpenseControllerTest struct {
	suite.Suite
	router    *gin.Engine
	expenseUC *usecase_mock.ExpensesUsecaseMock
	am        *middleware.AuthMiddleware
}

func (e *ExpenseControllerTest) SetupTest() {
	e.expenseUC = new(usecase_mock.ExpensesUsecaseMock)
	e.am = new(middleware.AuthMiddleware)

	e.router = gin.Default()
	gin.SetMode(gin.TestMode)

	rg := e.router.Group("/api/v1")

	// Mock Middleware untuk set "user" context key agar ctx.MustGet("user") tidak panic
	rg.Use(func(c *gin.Context) {
		c.Set("user", "uuid-user-1")
		c.Next()
	})

	expenseC := NewExpenseController(e.expenseUC, rg, *e.am)
	rg.POST("/expenses", expenseC.createHandler)
	rg.GET("/expenses", expenseC.listHandler)
	rg.GET("/expenses/:id", expenseC.getHandler)
	rg.GET("/expenses/transaction/:type", expenseC.getByTransactionHandler)
}

func TestExpenseControllerSuite(t *testing.T) {
	suite.Run(t, new(ExpenseControllerTest))
}

// TestCreateExpenseHandler_Success
func (e *ExpenseControllerTest) TestCreateExpenseHandler_Success() {
	// Arrange
	e.expenseUC.On("RegisterNewExpense", mock.Anything).Return(entity.Expense{}, nil).Once()

	payload := entity.Expense{
		TransactionType: "CREDIT",
		Amount:          100000,
		Date:            time.Now(),
		Balance:         100000,
		Description:     "Salary",
		UserId:          "uuid-user-1",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	e.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/expenses", &buf)
	e.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusCreated, record.Code)
}

// TestCreateExpenseHandler_Failed
func (e *ExpenseControllerTest) TestCreateExpenseHandler_Failed() {
	// Arrange
	e.expenseUC.On("RegisterNewExpense", mock.Anything).Return(entity.Expense{}, fmt.Errorf("failed")).Once()

	payload := entity.Expense{
		TransactionType: "CREDIT",
		Amount:          100000,
		Date:            time.Now(),
		Balance:         100000,
		Description:     "Salary",
		UserId:          "uuid-user-1",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	e.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/expenses", &buf)
	e.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusInternalServerError, record.Code)
}

// TestListExpenseHandler_Success
func (e *ExpenseControllerTest) TestListExpenseHandler_Success() {
	// prepare
	e.expenseUC.On("FindAllExpense", 1, 10, "", "", "uuid-user-1").Return([]entity.Expense{
		{
			ID:              "uuid-expense-1",
			TransactionType: "CREDIT",
			Amount:          100000,
			Date:            time.Now(),
			Balance:         100000,
			Description:     "Salary",
			UserId:          "uuid-user-1",
		},
	}, model.Paging{}, nil).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusOK, record.Code)
}

// TestListExpenseHandler_Failed
func (e *ExpenseControllerTest) TestListExpenseHandler_Failed() {
	// prepare
	e.expenseUC.On("FindAllExpense", 1, 10, "", "", "uuid-user-1").Return([]entity.Expense{}, model.Paging{}, fmt.Errorf("failed")).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusInternalServerError, record.Code)
}

// TestGetExpenseHandler_Success
func (e *ExpenseControllerTest) TestGetExpenseHandler_Success() {
	// prepare
	e.expenseUC.On("FindExpenseByID", "uuid-expense-1").Return(entity.Expense{
		ID:              "uuid-expense-1",
		TransactionType: "CREDIT",
		Amount:          100000,
		Date:            time.Now(),
		Balance:         100000,
		Description:     "Salary",
		UserId:          "uuid-user-1",
	}, nil).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses/uuid-expense-1", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusOK, record.Code)
}

// TestGetExpenseHandler_Failed
func (e *ExpenseControllerTest) TestGetExpenseHandler_Failed() {
	// prepare
	e.expenseUC.On("FindExpenseByID", "uuid-expense-1").Return(entity.Expense{}, fmt.Errorf("failed")).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses/uuid-expense-1", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusNotFound, record.Code)
}

// TestGetByTransactionHandler_Success
func (e *ExpenseControllerTest) TestGetByTransactionHandler_Success() {
	// prepare
	expenses := []entity.Expense{
		{
			ID:              "uuid-expense-1",
			TransactionType: "CREDIT",
			Amount:          100000,
			Date:            time.Now(),
			Balance:         100000,
			Description:     "Salary",
			UserId:          "uuid-user-1",
		},
	}

	e.expenseUC.On("FindExpenseByTransactionType", "CREDIT", "uuid-user-1").Return(expenses, nil).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses/transaction/CREDIT", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusOK, record.Code)
}

// TestGetByTransactionHandler_Failed
func (e *ExpenseControllerTest) TestGetByTransactionHandler_Failed() {
	// prepare
	e.expenseUC.On("FindExpenseByTransactionType", "CREDIT", "uuid-user-1").Return([]entity.Expense{}, fmt.Errorf("failed")).Once()

	req, err := http.NewRequest("GET", "/api/v1/expenses/transaction/CREDIT", nil)
	e.NoError(err)

	record := httptest.NewRecorder()

	// Act
	e.router.ServeHTTP(record, req)

	// Assert
	e.Equal(http.StatusNotFound, record.Code)
}