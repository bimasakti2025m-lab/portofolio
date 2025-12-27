// TODO :
// 1. Mendeklarasikan nama package pada file book_repositori
// 2. Mendeklarasikan struct bernama bookRepositori
// 3. Mendeklarasikan interface bernama BookRepositori
// 4. Membuat method dari interface yang telah dibuat
// 5. Mendeklarasikan constructor bernama newBookRepositori

package repositori

import (
	"database/sql"
	"simple-clean-architecture/model"
)

type bookRepositori struct {
	db *sql.DB
}

type BookRepositori interface {
	CreateNewBook(book model.Book) (model.Book, error)
	GetAllBook() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) (model.Book, error)
	DeleteBook(id int) error
}

func (b *bookRepositori) CreateNewBook(book model.Book) (model.Book, error) {
	var bookId int

	err := b.db.QueryRow("INSERT INTO mst_book(title, author, release_year, pages) VALUES($1, $2, $3, $4) RETURNING id", book.Title, book.Author, book.ReleaseYear, book.Pages).Scan(&bookId) 

	if err != nil {
		return model.Book{}, err
	}
	book.Id = bookId

	return book, nil
}

func (b *bookRepositori) GetAllBook() ([]model.Book, error) {
	var books []model.Book

	rows, err := b.db.Query("SELECT * FROM mst_book")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var book model.Book

		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *bookRepositori) GetBookById(id int) (model.Book, error) {
	var book model.Book

	err := b.db.QueryRow("SELECT * FROM mst_book WHERE id = $1", id).Scan(&book.Id, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)

	if err != nil {
		return model.Book{}, err
	}

	return book, nil
}

func (b *bookRepositori) UpdateBook(book *model.Book) (model.Book, error) {
	_, err := b.db.Exec("UPDATE mst_book SET title = $1, author = $2, release_year = $3, pages = $4 WHERE id = $5", book.Title, book.Author, book.ReleaseYear, book.Pages, book.Id)

	if err != nil {
		return model.Book{}, err
	}

	return *book, nil
}

func (b *bookRepositori) DeleteBook(id int) error {
	_, err := b.db.Exec("DELETE FROM mst_book WHERE id = $1", id)

	if err != nil {
		return err
	}

	return nil

}

func NewBookRepositori(db *sql.DB) BookRepositori {
	return &bookRepositori{db:db}
}
