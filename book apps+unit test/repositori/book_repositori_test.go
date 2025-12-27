package repositori

import (
	"regexp"
	"simple-clean-architecture/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	book := model.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ReleaseYear: 2023,
		Pages:       100,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_book(title, author, release_year, pages) VALUES($1, $2, $3, $4) RETURNING id")).
		WithArgs(book.Title, book.Author, book.ReleaseYear, book.Pages).
		WillReturnRows(rows)

	createdBook, err := repo.CreateNewBook(book)
	assert.NoError(t, err)
	assert.Equal(t, 1, createdBook.Id)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNewBook_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	book := model.Book{Title: "Test Book", Author: "Test Author", ReleaseYear: 2023, Pages: 100}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_book(title, author, release_year, pages) VALUES($1, $2, $3, $4) RETURNING id")).
		WithArgs(book.Title, book.Author, book.ReleaseYear, book.Pages).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.CreateNewBook(book)
	assert.Error(t, err)
}

func TestGetAllBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "release_year", "pages"}).
		AddRow(1, "Book 1", "Author 1", 2020, 150).
		AddRow(2, "Book 2", "Author 2", 2021, 200)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_book")).WillReturnRows(rows)

	books, err := repo.GetAllBook()
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, 1, books[0].Id)
	assert.Equal(t, "Book 1", books[0].Title)
	assert.Equal(t, "Author 1", books[0].Author)
	assert.Equal(t, 2020, books[0].ReleaseYear)
	assert.Equal(t, 150, books[0].Pages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllBook_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_book")).WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetAllBook()
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "release_year", "pages"}).
		AddRow(1, "Book 1", "Author 1", 2020, 150)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_book WHERE id = $1")).WithArgs(1).WillReturnRows(rows)

	book, err := repo.GetBookById(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, book.Id)
	assert.Equal(t, "Book 1", book.Title)
	assert.Equal(t, "Author 1", book.Author)
	assert.Equal(t, 2020, book.ReleaseYear)
	assert.Equal(t, 150, book.Pages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookById_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM mst_book WHERE id = $1")).WithArgs(1).WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetBookById(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE mst_book SET title = $1, author = $2, release_year = $3, pages = $4 WHERE id = $5")).
		WithArgs("Updated Book", "Updated Author", 2021, 200, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedBook, err := repo.UpdateBook(&model.Book{Id: 1, Title: "Updated Book", Author: "Updated Author", ReleaseYear: 2021, Pages: 200})
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedBook.Id)
	assert.Equal(t, "Updated Book", updatedBook.Title)
	assert.Equal(t, "Updated Author", updatedBook.Author)
	assert.Equal(t, 2021, updatedBook.ReleaseYear)
	assert.Equal(t, 200, updatedBook.Pages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE mst_book SET title = $1, author = $2, release_year = $3, pages = $4 WHERE id = $5")).
		WithArgs("Updated Book", "Updated Author", 2021, 200, 1).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.UpdateBook(&model.Book{Id: 1, Title: "Updated Book", Author: "Updated Author", ReleaseYear: 2021, Pages: 200})
	assert.Error(t, err)
}

func TestDeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_book WHERE id = $1")).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteBook(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBook_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBookRepositori(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM mst_book WHERE id = $1")).WithArgs(1).WillReturnError(sqlmock.ErrCancelled)

	err = repo.DeleteBook(1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
