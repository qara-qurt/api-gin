package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qara-qurt/api-gin/pkg/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId int, listId int) (model.TodoList, error)
	Update(userId int, listId int, input model.UpdateListInput) error
	Delete(userId int, listId int) error
}

type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
