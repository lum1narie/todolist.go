package service

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

const (
	salt = "f3f66d75-51fe-455d-a6f6-379902e4efc3"
)

func hash(pw string) []byte {
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(pw))
	return h.Sum(nil)
}

func RegisterUser(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	switch {
	case username == "":
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":    "Register user",
				"Error":    "Usernane is not provided",
				"Username": username,
			})
	case password == "":
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":    "Register user",
				"Error":    "Password is not provided",
				"Password": password,
			})
	}

	// DB 接続
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 重複チェック
	var duplicate int
	err = db.Get(&duplicate,
		"SELECT COUNT(*) FROM users WHERE name=?", username)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	if duplicate > 0 {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title": "Register user",
				"Error": fmt.Sprintf(
					"Username '%s' is already taken", username),
				"Username": username,
				"Password": password,
			})
		return
	}

	// DB への保存
	result, err := db.Exec(
		"INSERT INTO users(name, password) VALUES (?, ?)",
		username, hash(password))
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// 保存状態の確認
	id, _ := result.LastInsertId()
	var user database.User
	err = db.Get(&user, "SELECT id, name, password FROM users WHERE id = ?", id)
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func NewUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}
