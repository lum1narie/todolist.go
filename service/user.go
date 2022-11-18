package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

var (
	badPasswordLengthRe  = regexp.MustCompile(`^.{0,7}$`)
	badPasswordOnlyNumRe = regexp.MustCompile(`^[0-9]*$`)
)

const (
	salt    = "f3f66d7551fe455da6f6379902e4efc3"
	userkey = "user"
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

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// ユーザの取得
	var user database.User
	err = db.Get(&user,
		"SELECT id, name, password FROM users WHERE name = ?", username)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "login.html",
			gin.H{"Title": "Login", "Username": username,
				"Error": "User or Password is wrong"})
		return
	}

	// パスワードの照合
	if hex.EncodeToString(user.Password) != hex.EncodeToString(hash(password)) {
		ctx.HTML(http.StatusBadRequest, "login.html",
			gin.H{"Title": "Login", "Username": username,
				"Error": "User or Password is wrong"})
		return
	}

	// セッションの保存
	session := sessions.Default(ctx)
	session.Set(userkey, user.ID)
	session.Save()

	ctx.Redirect(http.StatusFound, "/list")
}

func Logout(ctx *gin.Context) {
    session := sessions.Default(ctx)
    session.Clear()
    session.Options(sessions.Options{MaxAge: -1})
    session.Save()
    ctx.Redirect(http.StatusFound, "/")
}

func LoginForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
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
