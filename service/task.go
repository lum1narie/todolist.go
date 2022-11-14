package service

import (
	"fmt"
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
	const maxTitleLen = 30
	const maxDescLen = 40
	const format = "2006-01-02 15:04:05"

	// truncate task strings
	title := task.Title
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}

	description := task.Description
	if len(description) > maxDescLen {
		description = description[:maxDescLen-3] + "..."
	}

	return &TaskViewInList{
		ID:          task.ID,
		Title:       title,
		Description: description,
		CreatedAt:   task.CreatedAt.Format(format),
		IsDone:      task.IsDone,
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
	// Use DB#Select for multiple entries
	err = db.Select(&tasks, "SELECT * FROM tasks") 
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
	ctx.HTML(http.StatusOK, "task_list.html",
		gin.H{"Title": "Task list", "Tasks": taskViews})
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
	// Use DB#Get for one entry
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task.html", task)
}

func NewTaskForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_new_task.html",
		gin.H{"Title": "Task registration"})
}

func RegisterTask(ctx *gin.Context) {
	// Get task title
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "No title is given")(ctx)
		return
	}

	// Get task description
	description, exist := ctx.GetPostForm("description")
	if !exist {
		description = ""
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Create new data with given title on DB
	result, err := db.Exec(
		"INSERT INTO tasks (title, description) VALUES (?, ?)",
		title, description)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	path := "/list" // デフォルトではタスク一覧ページへ戻る
	if id, err := result.LastInsertId(); err == nil {
		path = fmt.Sprintf("/task/%d", id) // 正常にIDを取得できた場合は /task/<id> へ戻る
	}
	ctx.Redirect(http.StatusFound, path)
}

func EditTaskForm(ctx *gin.Context) {
	// ID の取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	// Get target task
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Render edit form
	ctx.HTML(http.StatusOK, "form_edit_task.html",
		gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task})
}

func UpdateTask(ctx *gin.Context) {
	// ID の取得
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get task title
	title, exist := ctx.GetPostForm("title")
	if !exist {
		Error(http.StatusBadRequest, "No title is given")(ctx)
		return
	}

	// Get task description
	description, exist := ctx.GetPostForm("description")
	if !exist {
		description = ""
	}

	// Get task status
	is_done_raw, exist := ctx.GetPostForm("is_done")
	if !exist {
		Error(http.StatusBadRequest, "No is_done is given")(ctx)
		return
	}
	is_done, err := strconv.ParseBool(is_done_raw)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Create new data with given title on DB
	_, err = db.Exec(
		"UPDATE tasks SET title = ?, description = ?, is_done = ? WHERE id = ?",
		title, description, is_done, id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	path := fmt.Sprintf("/task/%d", id) 
	ctx.Redirect(http.StatusFound, path)
}
