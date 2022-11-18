package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"todolist.go/db"
	"todolist.go/service"
)

const (
	port         = 8000
	cookieSecret = "991a6ccaf0fc494e80263df738dcc601"
	sessionName  = "user-session"
)

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// prepare session
	store := cookie.NewStore([]byte(cookieSecret))
	engine.Use(sessions.Sessions(sessionName, store))

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)

	// ユーザ登録
	engine.GET("/user/new", service.NewUserForm)
	engine.POST("/user/new", service.RegisterUser)

	// ログイン
	engine.GET("/login", service.LoginForm)
	engine.POST("/login", service.Login)
	engine.GET("/logout", service.Logout)

	engine.GET("/list", service.LoginCheck, service.TaskList)

	taskGroup := engine.Group("/task")
	taskGroup.Use(service.LoginCheck)
	{
		taskGroup.GET("/:id", service.ShowTask) // ":id" is a parameter
		// タスクの新規登録
		taskGroup.GET("/new", service.NewTaskForm)
		taskGroup.POST("/new", service.RegisterTask)
		// 既存タスクの編集
		taskGroup.GET("/edit/:id", service.EditTaskForm)
		taskGroup.POST("/edit/:id", service.UpdateTask)
		// 既存タスクの削除
		taskGroup.GET("/delete/:id", service.DeleteTask)
	}

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
