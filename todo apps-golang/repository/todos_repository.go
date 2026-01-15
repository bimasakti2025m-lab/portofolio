package repository

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/golang-jwt/model"
)

// Mendeklarasikan interface
type TodosRepository interface {
	Create(model.Todo) (model.Todo, error)
	List() ([]model.Todo, error)
	Get(id string) (model.Todo, error)
	Update(id string, todo model.Todo) (model.Todo, error)
	Delete(id string) error
}

// Mendeklarasikan struct
type todosRepository struct {
	db *sql.DB
}

// Implementasi dari interface

func (t *todosRepository) Create(todo model.Todo) (model.Todo, error) {
	// Menambahkan user_id dan penanganan error
	err := t.db.QueryRow("INSERT INTO mst_todos (title, description, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at", todo.Title, todo.Description, todo.UserID).Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to create todo: %w", err)
	}
	return todo, nil
}

func (t *todosRepository) List() ([]model.Todo, error) {
	rows, err := t.db.Query("SELECT id, title, description, completed, user_id, created_at, updated_at FROM mst_todos")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo row: %w", err)
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (t *todosRepository) Get(id string) (model.Todo, error) {
	var todo model.Todo
	err := t.db.QueryRow("SELECT id, title, description, completed, user_id, created_at, updated_at FROM mst_todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Todo{}, fmt.Errorf("todo with id %s not found", id)
		}
		return model.Todo{}, fmt.Errorf("failed to get todo by id: %w", err)
	}
	return todo, nil
}

func (t *todosRepository) Update(id string, todo model.Todo) (model.Todo, error) {
	// Menggunakan COALESCE untuk memperbarui hanya field yang tidak kosong/nol
	// Ini adalah pendekatan sederhana, cara yang lebih baik adalah membangun query secara dinamis.
	err := t.db.QueryRow(`
		UPDATE mst_todos SET 
			title = COALESCE(NULLIF($2, ''), title), 
			description = COALESCE(NULLIF($3, ''), description), 
			completed = $4,
			updated_at = NOW()
		WHERE id = $1 AND user_id = $5 RETURNING id, title, description, completed, user_id, created_at, updated_at`,
		id, todo.Title, todo.Description, todo.Completed).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return model.Todo{}, fmt.Errorf("failed to update todo: %w", err)
	}
	return todo, nil
}

func (t *todosRepository) Delete(id string) error {
	result, err := t.db.Exec("DELETE FROM mst_todos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}
	return nil
}

// Mendeklarasikan konstruktor
func NewTodosRepository(db *sql.DB) TodosRepository {
	return &todosRepository{db: db}
}
