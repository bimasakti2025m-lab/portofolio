package usecase

import (
	"simple-clean-architecture/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateNewBook(book model.Book) (model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookRepository) GetAllBook() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookRepository) GetBookById(id int) (model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookRepository) UpdateBook(book *model.Book) (model.Book, error) {
	args := m.Called(book)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookRepository) DeleteBook(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookUsecase(t *testing.T) {
	repo := new(MockBookRepository)
	repo.On("CreateNewBook", mock.Anything).Return(model.Book{}, nil)
	repo.On("GetAllBook").Return([]model.Book{}, nil)
	repo.On("GetBookById", mock.Anything).Return(model.Book{}, nil)
	repo.On("UpdateBook", mock.Anything).Return(model.Book{}, nil)
	repo.On("DeleteBook", mock.Anything).Return(nil)

	usecase := NewBookUsecase(repo)

	book := model.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ReleaseYear: 2023,
		Pages:       100,
	}

	createdBook, err := usecase.CreateNewBook(book)
	assert.NoError(t, err)
	assert.NotNil(t, createdBook)

	books, err := usecase.GetAllBook()
	assert.NoError(t, err)
	assert.NotNil(t, books)

	bookById, err := usecase.GetBookById(1)
	assert.NoError(t, err)
	assert.NotNil(t, bookById)

	// validation if Pages is 100
	assert.Equal(t, 100, book.Pages)

	updatedBook, err := usecase.UpdateBook(&book)
	assert.NoError(t, err)
	assert.NotNil(t, updatedBook)

	err = usecase.DeleteBook(1)
	assert.NoError(t, err)
}