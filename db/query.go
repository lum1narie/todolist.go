package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SearchBool int

const (
	SearchTrue SearchBool = iota
	SearchFalse
	SearchBoth
)

func GetUserById(db *sqlx.DB, id uint64) (User, error) {
	var user User
	err := db.Get(&user,
		"SELECT id, name, password FROM users WHERE id = ?", id)
	return user, err
}

func GetUserByName(db *sqlx.DB, name string) (User, error) {
	var user User
	err := db.Get(&user,
		"SELECT id, name, password FROM users WHERE name = ?", name)
	return user, err
}

func IsUserWithNameExist(db *sqlx.DB, name string) (bool, error) {
	var count int
	err := db.Get(&count,
		"SELECT COUNT(*) FROM users WHERE name=?", name)
	return count > 0, err
}

func AddUser(db *sqlx.DB, name string, password []byte) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO users(name, password) VALUES (?, ?)",
		name, password)
}

func GetTaskById(db *sqlx.DB, id uint64) (Task, error) {
	var task Task
	// Use DB#Get for one entry
	err := db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	return task, err
}

func GetTasksByUser(db *sqlx.DB,
	userID uint64, name string, status SearchBool) ([]Task, error) {
	var tasks []Task
	var err error

	query := "SELECT id, title, description, created_at, is_done" +
		" FROM tasks INNER JOIN ownership ON task_id = tasks.id" +
		" WHERE user_id = ?"
	switch status {
	case SearchTrue:
		query += " AND is_done = TRUE"
	case SearchFalse:
		query += " AND is_done = FALSE"
	}
	if name == "" {
		err = db.Select(&tasks, query, userID)
	} else {
		query += " AND title LIKE ?"
		err = db.Select(&tasks, query, userID, "%"+name+"%")
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func AddTask(db *sqlx.DB, title string, description string) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO tasks (title, description) VALUES (?, ?)",
		title, description)
}

func UpdateTaskById(db *sqlx.DB, id uint64,
	title string, description string, is_done bool) (sql.Result, error) {
	return db.Exec(
		"UPDATE tasks SET title = ?, description = ?, is_done = ? WHERE id = ?",
		title, description, is_done, id)
}

func DeleteTaskById(db *sqlx.DB, id uint64) (sql.Result, error) {
	return db.Exec("DELETE FROM tasks WHERE id=?", id)
}
