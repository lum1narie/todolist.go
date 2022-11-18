package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func checkIsOwnedTask(ctx *gin.Context, db *sqlx.DB, taskID uint64) (bool, error) {
	userID, _ := sessions.Default(ctx).Get(userkey).(uint64)
	return database.IsTaskBelongsToUser(db, taskID, userID)
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Get query parameter
	kw := ctx.Query("kw")
	var status database.SearchBool
	statusRaw := ctx.Query("status")
	switch statusRaw {
	case "finished":
		status = database.SearchTrue
	case "unfinished":
		status = database.SearchFalse
	default:
		status = database.SearchBoth
	}

	// Get tasks in DB
	userID, _ := sessions.Default(ctx).Get(userkey).(uint64)
	tasks, err := database.GetTasksByUser(db, userID, kw, status)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// truncate task description
	var taskViews []TaskViewInList
	for _, task := range tasks {
		taskViews = append(taskViews, *viewInListFromTask(&task))
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html",
		gin.H{"Title": "Task list", "Tasks": taskViews, "Kw": kw,
			"Status": statusRaw})
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

	// reject if user is not owner
	isOwner, err := checkIsOwnedTask(ctx, db, uint64(id))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	if !isOwner {
		NotFound()(ctx)
		return
	}

	// Get a task with given ID
	task, err := database.GetTaskById(db, uint64(id))
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
	userID, _ := sessions.Default(ctx).Get(userkey).(uint64)

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
	taskID, err := database.AddTaskWithUser(db, title, description, userID)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Render status
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/task/%d", taskID))
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

	// reject if user is not owner
	isOwner, err := checkIsOwnedTask(ctx, db, uint64(id))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	if !isOwner {
		NotFound()(ctx)
		return
	}

	// Get target task
	task, err := database.GetTaskById(db, uint64(id))
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

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// reject if user is not owner
	isOwner, err := checkIsOwnedTask(ctx, db, uint64(id))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	if !isOwner {
		NotFound()(ctx)
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
	isDoneRaw, exist := ctx.GetPostForm("is_done")
	if !exist {
		Error(http.StatusBadRequest, "No is_done is given")(ctx)
		return
	}
	isDone, err := strconv.ParseBool(isDoneRaw)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Create new data with given title on DB
	_, err = database.UpdateTaskById(db, uint64(id), title, description, isDone)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	path := fmt.Sprintf("/task/%d", id)
	ctx.Redirect(http.StatusFound, path)
}

func DeleteTask(ctx *gin.Context) {
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

	// reject if user is not owner
	isOwner, err := checkIsOwnedTask(ctx, db, uint64(id))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	if !isOwner {
		NotFound()(ctx)
		return
	}

	// Delete the task from DB
	_, err = database.DeleteTaskById(db, uint64(id))
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	// Redirect to /list
	ctx.Redirect(http.StatusFound, "/list")
}
