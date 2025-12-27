package usecase

import (
	"errors"
	"testing"
	"time"

	"enigmacamp.com/livecode-catatan-keuangan/entity"
	"enigmacamp.com/livecode-catatan-keuangan/mock/usecase_mock"
	"enigmacamp.com/livecode-catatan-keuangan/shared/model"
	"github.com/stretchr/testify/suite"
)

type ExpensesUCSuite struct {
	suite.Suite
	expenseRepo *usecase_mock.ExpensesUsecaseMock
	expenseUC   ExpenseUseCase
}

func (e *ExpensesUCSuite) SetupTest() {
	e.expenseRepo = new(usecase_mock.ExpensesUsecaseMock)
	e.expenseUC = NewExpenseUseCase(e.expenseRepo)
}
func TestExpensesUCSuite(t *testing.T) {
	suite.Run(t, new(ExpensesUCSuite))
}

func (e *ExpensesUCSuite) TestFindAllExpense_success() {
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
		{
			ID:              "uuid-expense-2",
			TransactionType: "DEBIT",
			Amount:          50000,
			Date:            time.Now(),
			Balance:         50000,
			Description:     "Food",
			UserId:          "uuid-user-1",
		},
	}

	e.expenseRepo.On("List", 1, 10, "", "", "uuid-user-1").Return(expenses, model.Paging{}, nil).Once()

	result, paging, err := e.expenseUC.FindAllExpense(1, 10, "", "", "uuid-user-1")
	e.Nil(err)
	e.Equal(expenses, result)
	e.Equal(model.Paging{}, paging)
}

func (e *ExpensesUCSuite) TestFindExpenseByID_success() {
	expense := entity.Expense{
		ID:              "uuid-expense-1",
		TransactionType: "CREDIT",
		Amount:          100000,
		Date:            time.Now(),
		Balance:         100000,
		Description:     "Salary",
		UserId:          "uuid-user-1",
	}

	e.expenseRepo.On("Get", "uuid-expense-1").Return(expense, nil).Once()

	result, err := e.expenseUC.FindExpenseByID("uuid-expense-1")
	e.Nil(err)
	e.Equal(expense, result)
}

func (e *ExpensesUCSuite) TestFindExpenseByTransactionType_success() {
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
		{
			ID:              "uuid-expense-3",
			TransactionType: "CREDIT",
			Amount:          200000,
			Date:            time.Now(),
			Balance:         300000,
			Description:     "Salary",
			UserId:          "uuid-user-1",
		},
	}

	e.expenseRepo.On("GetByTransaction", "CREDIT", "uuid-user-1").Return(expenses, nil).Once()

	result, err := e.expenseUC.FindExpenseByTransactionType("CREDIT", "uuid-user-1")
	e.Nil(err)
	e.Equal(expenses, result)
}

func (e *ExpensesUCSuite) TestRegisterNewExpense_success() {
	// prepare
	newExpense := entity.Expense{
		TransactionType: "CREDIT",
		Balance:         350000.0,
		Amount:          150000,
		Description:     "Freelance Project",
		UserId:          "uuid-user-1",
		Date:            time.Now(),
		UpdatedAt:       time.Now(),
	}

	// mocking
	e.expenseRepo.On("GetBalance", "uuid-user-1").Return(200000.0, nil).Once()
	e.expenseRepo.On("Create", newExpense).Return(newExpense, nil).Once()

	// execute
	result, err := e.expenseUC.RegisterNewExpense(newExpense)

	// assert
	e.Nil(err)
	e.NotNil(result)
}

func (e *ExpensesUCSuite) TestRegisterNewExpense_failed() {
	// prepare
	newExpense := entity.Expense{
		TransactionType: "CREDIT",
		Balance:         350000.0,
		Amount:          150000,
		Description:     "Freelance Project",
		UserId:          "uuid-user-1",
		Date:            time.Now(),
		UpdatedAt:       time.Now(),
	}

	// mocking
	e.expenseRepo.On("GetBalance", "uuid-user-1").Return(200000.0, nil).Once()
	e.expenseRepo.On("Create", newExpense).Return(entity.Expense{}, errors.New("error")).Once()

	// execute
	result, err := e.expenseUC.RegisterNewExpense(newExpense)

	// assert
	e.NotNil(err)
	e.Equal(entity.Expense{}, result)
}

func (e *ExpensesUCSuite) TestRegisterNewExpense_invalidTransactionType() {
	// prepare
	newExpense := entity.Expense{
		TransactionType: "INVALID",
		Balance:         350000.0,
		Amount:          150000,
		Description:     "Freelance Project",
		UserId:          "uuid-user-1",
		Date:            time.Now(),
		UpdatedAt:       time.Now(),
	}

	// execute
	result, err := e.expenseUC.RegisterNewExpense(newExpense)

	// assert
	e.NotNil(err)
	e.Equal(entity.Expense{}, result)
}

func (e *ExpensesUCSuite) TestRegisterNewExpense_invalidAmount() {
	// prepare
	newExpense := entity.Expense{
		TransactionType: "CREDIT",
		Balance:         50000.0,
		Amount:          -150000,
		Description:     "Freelance Project",
		UserId:          "uuid-user-1",
		Date:            time.Now(),
		UpdatedAt:       time.Now(),
	}

	// mocking
	e.expenseRepo.On("GetBalance", "uuid-user-1").Return(200000.0, nil).Once()
	e.expenseRepo.On("Create", newExpense).Return(entity.Expense{}, errors.New("error")).Once()

	// execute
	result, err := e.expenseUC.RegisterNewExpense(newExpense)

	// assert
	e.NotNil(err)
	e.Equal(entity.Expense{}, result)
}

func (e *ExpensesUCSuite) TestRegisterNewExpense_invalidBalance() {
	// prepare
	newExpense := entity.Expense{
		TransactionType: "CREDIT",
		Balance:         350000.0,
		Amount:          150000,
		Description:     "Freelance Project",
		UserId:          "uuid-user-1",
		Date:            time.Now(),
		UpdatedAt:       time.Now(),
	}

	// mocking
	e.expenseRepo.On("GetBalance", "uuid-user-1").Return(200000.0, nil).Once()
	e.expenseRepo.On("Create", newExpense).Return(entity.Expense{}, errors.New("error")).Once()

	// execute
	result, err := e.expenseUC.RegisterNewExpense(newExpense)

	// assert
	e.NotNil(err)
	e.Equal(entity.Expense{}, result)
}
