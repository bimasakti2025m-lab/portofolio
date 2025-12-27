package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"simple-clean-architecture/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookUsecase struct {
	mock.Mock
}

func (m *MockBookUsecase) CreateNewBook(book model.Book) (model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUsecase) GetAllBook() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookUsecase) GetBookById(id int) (model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUsecase) UpdateBook(book *model.Book) (model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUsecase) DeleteBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookController_CreateNewBook(t *testing.T) {
	mockUsecase := new(MockBookUsecase)
	router := gin.Default()
	rg := router.Group("/api/v1")
	NewBookController(mockUsecase, rg).Route()

	// Happy Path
	book := model.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ReleaseYear: 2023,
		Pages:       100,
	}

	// Setup mock
	// Kita gunakan mock.MatchedBy untuk memastikan argumen yang masuk sesuai
	mockUsecase.On("CreateNewBook", mock.MatchedBy(func(b model.Book) bool {
		return b.Title == book.Title
	})).Return(model.Book{Id: 1, Title: "Test Book", Author: "Test Author", ReleaseYear: 2023, Pages: 100}, nil).Once()

	body, err := json.Marshal(book)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUsecase.AssertExpectations(t)

	// Sad Path: Invalid JSON
	req, err = http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer([]byte(`{"title": "bad json"`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookController_GetAllBook(t *testing.T) {
	mockUsecase := new(MockBookUsecase)
	router := gin.Default()
	rg := router.Group("/api/v1")
	NewBookController(mockUsecase, rg).Route()

	// Happy Path
	expectedBooks := []model.Book{{Id: 1, Title: "Book 1"}}
	mockUsecase.On("GetAllBook").Return(expectedBooks, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualBooks []model.Book
	err = json.Unmarshal(w.Body.Bytes(), &actualBooks)
	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, actualBooks)
	mockUsecase.AssertExpectations(t)
}

func TestBookController_GetBookById(t *testing.T) {
	mockUsecase := new(MockBookUsecase)
	router := gin.Default()
	rg := router.Group("/api/v1")
	NewBookController(mockUsecase, rg).Route()

	// Happy Path
	expectedBook := model.Book{Id: 1, Title: "Book 1"}
	mockUsecase.On("GetBookById", 1).Return(expectedBook, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/books/1", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualBook model.Book
	err = json.Unmarshal(w.Body.Bytes(), &actualBook)
	assert.NoError(t, err)
	assert.Equal(t, expectedBook, actualBook)
	mockUsecase.AssertExpectations(t)

	// Sad Path: Usecase returns error
	mockUsecase.On("GetBookById", 2).Return(model.Book{}, errors.New("not found")).Once()

	req, err = http.NewRequest(http.MethodGet, "/api/v1/books/2", nil)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBookController_UpdateBook(t *testing.T) {
	mockUsecase := new(MockBookUsecase)
	router := gin.Default()
	rg := router.Group("/api/v1")
	NewBookController(mockUsecase, rg).Route()

	// Happy Path
	book := model.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ReleaseYear: 2023,
		Pages:       100,
	}

	mockUsecase.On("UpdateBook", &model.Book{Id: 1, Title: "Test Book", Author: "Test Author", ReleaseYear: 2023, Pages: 100}).Return(model.Book{Id: 1, Title: "Test Book", Author: "Test Author", ReleaseYear: 2023, Pages: 100}, nil).Once()

	body, err := json.Marshal(book)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/api/v1/books/1", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestBookController_DeleteBook(t *testing.T) {
	mockUsecase := new(MockBookUsecase)
	router := gin.Default()
	rg := router.Group("/api/v1")
	NewBookController(mockUsecase, rg).Route()

	// Happy Path
	mockUsecase.On("DeleteBook", 1).Return(nil).Once()

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/books/1", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Book deleted successfully", response["message"])
	mockUsecase.AssertExpectations(t)
}
