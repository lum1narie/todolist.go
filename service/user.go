package service

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

var (
	badPasswordLengthRe  = regexp.MustCompile(`^.{0,7}$`)
	badPasswordOnlyNumRe = regexp.MustCompile(`^[0-9]*$`)
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

func validatePassword(pw string) (bool, string) {
	if badPasswordLengthRe.MatchString(pw) {
		return false, "Password is too short"
	}
	if badPasswordOnlyNumRe.MatchString(pw) {
		return false, "Password only contains number"
	}

	return true, ""
}

func RegisterUser(ctx *gin.Context) {
	// フォームデータの受け取り
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	passwordRe := ctx.PostForm("password-re")

	if username == "" {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":      "Register user",
				"Error":      "Usernane is not provided",
				"Password":   password,
				"PasswordRe": passwordRe,
			})
		return
	}
	if password == "" {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":    "Register user",
				"Error":    "Password is not provided",
				"Username": username,
			})
		return
	}
	if password != passwordRe {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":    "Register user",
				"Error":    "Retype Password doesn't match",
				"Username": username,
			})
		return
	}

	if result, msg := validatePassword(password); !result {
		ctx.HTML(http.StatusBadRequest, "new_user_form.html",
			gin.H{
				"Title":    "Register user",
				"Error":    msg,
				"Username": username,
			})
		return
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
				"Username":   username,
				"Password":   password,
				"PasswordRe": passwordRe,
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
