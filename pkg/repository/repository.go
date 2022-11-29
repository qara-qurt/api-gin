package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qara-qurt/api-gin/pkg/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
}

type TodoList interface{}

type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
