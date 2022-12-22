package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/qara-qurt/api-gin/pkg/model"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) Create(listId int, item model.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title,description) VALUES($1, $2) RETURNING id", todoItemTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (item_id, list_id) VALUES($1, $2)", listsItemTable)
	if _, err := tx.Exec(createListItemQuery, itemId, listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]model.TodoItem, error) {

	var items []model.TodoItem
	getAllItemsQuery :=
		fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done
						FROM %s ti 
						JOIN %s li ON ti.id = li.item_id
						JOIN %s ul ON ul.list_id = li.list_id
						WHERE li.list_id = $1 AND ul.user_id = $2`,
			todoItemTable, listsItemTable, usersListTable)

	if err := r.db.Select(&items, getAllItemsQuery, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId int, itemId int) (model.TodoItem, error) {
	var item model.TodoItem

	getItemByIdQuery :=
		fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
						FROM %s ti 
						INNER JOIN %s li on li.item_id = ti.id
						INNER JOIN %s ul on ul.list_id = li.list_id 
						WHERE ti.id = $1 AND ul.user_id = $2`,
			todoItemTable, listsItemTable, usersListTable)

	err := r.db.Get(&item, getItemByIdQuery, itemId, userId)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Update(userId int, itemId int, input model.UpdateItemInput) error {
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
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	// title = $1
	// description = $2
	// done = $3
	setQuery := strings.Join(setValues, ", ")

	updateQuery :=
		fmt.Sprintf(`UPDATE %s ti 
						SET %s 
						FROM %s li, %s ul
						WHERE ti.id = li.item_id 
						AND li.list_id = ul.list_id 
						AND ul.user_id = $%d 
						AND ti.id = $%d`,
			todoItemTable, setQuery, listsItemTable, usersListTable, argId, argId+1)

	args = append(args, userId, itemId)

	logrus.Debugf("updateQuery : %s", updateQuery)
	logrus.Debugf("args : %s", args)

	_, err := r.db.Exec(updateQuery, args...)

	return err
}

func (r *TodoItemPostgres) Delete(userId int, itemId int) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s ti 
								USING %s li, %s ul  
								WHERE ti.id = li.item_id 
								AND li.list_id = ul.list_id 
								AND ul.user_id = $1 
								AND ti.id = $2`,
		todoItemTable, listsItemTable, usersListTable)

	_, err := r.db.Exec(deleteQuery, userId, itemId)

	return err
}
