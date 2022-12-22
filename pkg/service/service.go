package service

import (
	"github.com/qara-qurt/api-gin/pkg/model"
	"github.com/qara-qurt/api-gin/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId int, listId int) (model.TodoList, error)
	Update(userId int, listId int, input model.UpdateListInput) error
	Delete(userId int, listId int) error
}

type TodoItem interface {
	Create(userId int, listId int, item model.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]model.TodoItem, error)
	GetById(userId int, itemId int) (model.TodoItem, error)
	Update(userId int, itemId int, input model.UpdateItemInput) error
	Delete(userId int, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
