package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pudongping/goblog/app/http/controllers"
	"github.com/pudongping/goblog/app/http/middlewares"
)

// RegisterWebRoutes 注册网页相关路由
func RegisterWebRoutes(r *mux.Router) {

	// 静态页面
	pc := new(controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	// 文章相关页面
	ac := new(controllers.ArticlesController)
	// 首页
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	// 文章详情
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	// 文章列表
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	// 文章创建显示
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	// 文章创建数据提交
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	// 更新文章显示
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	// 更新文章数据提交
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	// 删除文章
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	// 用户认证
	auc := new(controllers.AuthController)
	// 注册页面
	r.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	// 处理注册逻辑
	r.HandleFunc("/auth/do-register", auc.DoRegister).Methods("POST").Name("auth.doregister")
	// 登录显示页面
	r.HandleFunc("/auth/login", auc.Login).Methods("GET").Name("auth.login")
	// 处理登录逻辑
	r.HandleFunc("/auth/dologin", auc.DoLogin).Methods("POST").Name("auth.dologin")

	// 设置静态资源路由
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件：强制内容类型为 HTML
	// r.Use(middlewares.ForceHTML)
	// 中间件：开始会话
	r.Use(middlewares.StartSession)
}
