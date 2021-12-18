package controllers

import (
	"fmt"
	"net/http"

	"github.com/pudongping/goblog/app/models/user"
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

	// 1. 表单验证
	// 2. 验证通过，则入库并跳转到首页
	_user := user.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	_user.Create()

	if _user.ID > 0 {
		fmt.Fprint(w, "插入成功，ID为"+_user.GetStringID())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "创建用户失败，请联系管理员")
	}

	// 3. 表单不通过，则重新显示表单
}
