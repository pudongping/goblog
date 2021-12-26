package bootstrap

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pudongping/goblog/pkg/route"
	"github.com/pudongping/goblog/routes"
)

// SetupRoute 路由初始化
func SetupRoute(staticFS embed.FS) *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetRoute(router)

	// 静态资源打包
	// 相当于去除了 public 前缀
	// 原来的路径：http://localhost:3000/public/css/app.css
	// 我们修改后的访问路径为：http://localhost:3000/css/app.css
	sub, _ := fs.Sub(staticFS, "public")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))

	return router
}