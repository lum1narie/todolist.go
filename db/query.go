package db

import (
	"github.com/jmoiron/sqlx"
)

type SearchBool int

const (
	SearchTrue SearchBool = iota
	SearchFalse
	SearchBoth
)

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
