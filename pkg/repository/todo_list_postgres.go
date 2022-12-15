package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/qara-qurt/api-gin/pkg/model"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list model.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES($1, $2) RETURNING id", todoListTable)
	row := r.db.QueryRow(createListQuery, list.Title, list.Description)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListTable)
	_, err = r.db.Exec(createUserListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]model.TodoList, error) {
	var lists []model.TodoList

	getAllListsQuery := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListTable,
		usersListTable,
	)

	err := r.db.Select(&lists, getAllListsQuery, userId)
	return lists, err

}

func (r *TodoListPostgres) GetById(userId, listId int) (model.TodoList, error) {
	var list model.TodoList

	getListByIdQuery := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListTable,
		usersListTable,
	)

	err := r.db.Get(&list, getListByIdQuery, userId, listId)
	return list, err

}
