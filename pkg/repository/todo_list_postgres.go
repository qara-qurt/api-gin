package repository

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

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

func (r *TodoListPostgres) Update(userId int, listId int, input model.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateQuery := fmt.Sprintf(`UPDATE %s tl 
								SET %s FROM %s ul 
								WHERE tl.id = ul.list_id 
								AND ul.list_id = $%d 
								AND ul.user_id = $%d`,
		todoListTable, setQuery, usersListTable, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updateQuery : %s", updateQuery)
	logrus.Debugf("args : %s", args)

	_, err := r.db.Exec(updateQuery, args...)

	return err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {

	deleteQuery :=
		fmt.Sprintf(`DELETE FROM %s ti 
							USING %s li, %s ul 
							WHERE ti.id = li.item_id 
							AND li.list_id = ul.list_id 
							AND ul.user_id = $1 
							AND ti.id = $2`,
			todoItemTable, listsItemTable, usersListTable)

	_, err := r.db.Exec(deleteQuery, userId, listId)

	return err
}
