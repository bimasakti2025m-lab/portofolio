// TODO :
// 1. Mendeklarasikan nama package usecase
// 2. Mendeklarasikan struct bernama bookUsecase
// 3. Mendeklarasikan interface bernama BookUsecase
// 4. Membuat method dari interface yang telah dibuat
// 5. Mendeklarasikan constructor bernama NewBookUsecase

package usecase

import (
	"simple-clean-architecture/model"
	"simple-clean-architecture/repositori"

)

type bookUsecase struct {
	bookRepositori repositori.BookRepositori
}

type BookUsecase interface {
	CreateNewBook(book model.Book) (model.Book, error)
	GetAllBook() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) (model.Book, error)
	DeleteBook(id int) error
}

func (b *bookUsecase) CreateNewBook(book model.Book) (model.Book, error) {
	book, err := b.bookRepositori.CreateNewBook(book)

	if err != nil {
		return model.Book{}, err
	}

	return book, nil
}

func (b *bookUsecase) GetAllBook() ([]model.Book, error) {
	books, err := b.bookRepositori.GetAllBook()

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *bookUsecase) GetBookById(id int) (model.Book, error) {
	book, err := b.bookRepositori.GetBookById(id)

	if err != nil {
		return model.Book{}, err
	}

	return book, nil
}

func (b *bookUsecase) UpdateBook(book *model.Book) (model.Book, error) {
	updatedBook, err := b.bookRepositori.UpdateBook(book)

	if err != nil {
		return model.Book{}, err
	}

	return updatedBook, nil
}

func (b *bookUsecase) DeleteBook(id int) error {
	err := b.bookRepositori.DeleteBook(id)

	if err != nil {
		return err
	}

	return nil
}	

func NewBookUsecase(bookRepositori repositori.BookRepositori) BookUsecase {
	return &bookUsecase{bookRepositori: bookRepositori}
}
