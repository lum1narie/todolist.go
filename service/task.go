package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

type (
	TaskViewInList struct {
		ID          uint64
		Title       string
		Description string
		CreatedAt   string
		IsDone      bool
	}
)

func viewInListFromTask(task *database.Task) *TaskViewInList {
	const maxTitleLen = 20
	const maxDescLen = 30
	const format = "2006-01-02 15:04:05"

	// truncate task strings
	title := task.Title
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen - 3] + "..."
	}

	description := task.Description
	if len(description) > maxDescLen {
		description = description[:maxDescLen - 3] + "..."
	}

	return &TaskViewInList {
		ID: task.ID,
		Title: title,
		Description: description,
		CreatedAt: task.CreatedAt.Format(format),
		IsDone: task.IsDone,
	}
}


// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks") // Use DB#Select for multiple entries
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// truncate task description
	const maxDescLen = 30
	var taskViews []TaskViewInList
	for _, task := range tasks {
		taskViews = append(taskViews, *viewInListFromTask(&task))
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": taskViews})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task.html", task)
}
