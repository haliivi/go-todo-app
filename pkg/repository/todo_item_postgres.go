package repository

import (
	"fmt"
	"strings"

	"github.com/haliivi/go-todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, input todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println("3")
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemTable)
	var itemId int
	row := tx.QueryRow(createItemQuery, input.Title, input.Description)
	err = row.Scan(&itemId)
	if err != nil {
		fmt.Println("2")
		tx.Rollback()
		return 0, err
	}
	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		fmt.Println("1")
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN %s ul ON ul.list_id = li.list_id WHERE  li.list_id = $1 AND ul.user_id = $2", todoItemTable, listsItemsTable, usersListTable)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN %s ul ON ul.list_id = li.list_id WHERE  ti.id = $1 AND ul.user_id = $2", todoItemTable, listsItemsTable, usersListTable)
	err := r.db.Get(&item, query, itemId, itemId)
	return item, err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2", todoItemTable, listsItemsTable, usersListTable)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d", todoItemTable, setQuery, listsItemsTable, usersListTable, argId, argId+1)
	fmt.Println(query)
	args = append(args, userId, itemId)
	logrus.Debugf("updateQuery %s", query)
	logrus.Debugf("args %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
