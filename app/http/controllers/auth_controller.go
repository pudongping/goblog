package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thedevsaddam/govalidator"

	"github.com/pudongping/goblog/pkg/view"
)

type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// 验证器
type userForm struct {
	Name            string `valid:"name"`
	Email           string `valid:"email"`
	Password        string `valid:"password"`
	PasswordConfirm string `valid:"password_confirm"`
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 0. 初始化变量
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	passwordConfirm := r.PostFormValue("password_confirm")

	// 1. 初始化数据
	_user := userForm{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}

	// 2. 表单规则
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20"},
		"email":            []string{"required", "min:4", "max:30", "email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	// 3. 配置选项
	opts := govalidator.Options{
		Data:          &_user,
		Rules:         rules,
		TagIdentifier: "valid", // Struct 标签标识符
	}

	// 4. 开始验证
	errs := govalidator.New(opts).ValidateStruct()

	if len(errs) > 0 {
		// 4.1 有错误发生，打印错误数据
		// json.MarshalIndent 用来将 Go 对象格式成为 json 字符串，并加上合理的缩进
		data, _ := json.MarshalIndent(errs, "", "  ")
		fmt.Fprint(w, string(data))
	} else {
		// _user.Create()
		//
		// if _user.ID > 0 {
		// 	fmt.Fprint(w, "插入成功，ID为"+_user.GetStringID())
		// } else {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprint(w, "创建用户失败，请联系管理员")
		// }
	}

	// 3. 表单不通过，则重新显示表单
}
