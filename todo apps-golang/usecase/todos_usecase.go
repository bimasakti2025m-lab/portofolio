package usecase

import (
	"enigmacamp.com/golang-jwt/model"
	"enigmacamp.com/golang-jwt/repository"
)

// Mendeklarasikan interface
type TodosUseCase interface {
	Create(model.Todo) (model.Todo, error)
	List() ([]model.Todo, error)
	Get(id string) (model.Todo, error)
	Update(id string, todo model.Todo) (model.Todo, error)
	Delete(id string) error
}


// Mendeklarasikan struct
type todosUseCase struct {
	repo repository.TodosRepository
}

// Implementasi dari interface
func (t *todosUseCase) Create(todo model.Todo) (model.Todo, error) {
	return t.repo.Create(todo)
}

func (t *todosUseCase) List() ([]model.Todo, error) {
	return t.repo.List()
}

func (t *todosUseCase) Get(id string) (model.Todo, error) {
	return t.repo.Get(id)
}

func (t *todosUseCase) Update(id string, todo model.Todo) (model.Todo, error) {
	return t.repo.Update(id, todo)
}

func (t *todosUseCase) Delete(id string) error {
	return t.repo.Delete(id)
}


// Mendeklarasikan konstruktor
func NewTodosUseCase(repo repository.TodosRepository) TodosUseCase {
	return &todosUseCase{repo: repo}
}
