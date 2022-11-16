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

func GetTasks(db *sqlx.DB, name string, status SearchBool) ([]Task, error) {
	var tasks []Task
	var err error
	if name == "" {
		switch status {
		case SearchBoth:
			err = db.Select(&tasks, "SELECT * FROM tasks")
		case SearchTrue:
			err = db.Select(&tasks, "SELECT * FROM tasks WHERE is_done = TRUE")
		case SearchFalse:
			err = db.Select(&tasks, "SELECT * FROM tasks WHERE is_done = FALSE")
		}
	} else {
		switch status {
		case SearchBoth:
			err = db.Select(&tasks,
				"SELECT * FROM tasks WHERE title LIKE ?", "%"+name+"%")
		case SearchTrue:
			err = db.Select(&tasks,
				"SELECT * FROM tasks WHERE is_done = TRUE and title LIKE ?",
				"%"+name+"%")
		case SearchFalse:
			err = db.Select(&tasks,
				"SELECT * FROM tasks WHERE is_done = FALSE and title LIKE ?",
				"%"+name+"%")
		}
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
