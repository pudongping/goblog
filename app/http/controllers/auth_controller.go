package controllers

import (
	"fmt"
	"net/http"

	"github.com/pudongping/goblog/app/models/user"
	"github.com/pudongping/goblog/app/requests"
	"github.com/pudongping/goblog/pkg/auth"
	"github.com/pudongping/goblog/pkg/view"
)

type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 0. 初始化变量
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	passwordConfirm := r.PostFormValue("password_confirm")

	// 1. 初始化数据
	_user := user.User{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}

	// 2. 表单规则
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 3.1 有错误发生，打印错误数据
		// json.MarshalIndent 用来将 Go 对象格式成为 json 字符串，并加上合理的缩进
		// data, _ := json.MarshalIndent(errs, "", "  ")
		// fmt.Fprint(w, string(data))

		// 3. 表单不通过，重新显示表单
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")

	} else {
		// 4. 验证成功，创建数据
		_user.Create()

		if _user.ID > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "注册失败，请联系管理员")
		}
	}

}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登录表单提交
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 2. 尝试登录
	if err := auth.Attempt(email, password); err == nil {
		// 登录成功
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// 3. 登录失败，显示错误提示
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}
