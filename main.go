package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pudongping/goblog/app/http/middlewares"
	"github.com/pudongping/goblog/bootstrap"
	"github.com/pudongping/goblog/config"
	c "github.com/pudongping/goblog/pkg/config"
	"github.com/pudongping/goblog/pkg/logger"
)

var router *mux.Router

//var db *sql.DB

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {

	//database.Initialize()
	//db = database.DB

	// 初始化 gorm
	bootstrap.SetupDB()
	// 初始化路由绑定
	router = bootstrap.SetupRoute()

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)

	err := http.ListenAndServe(fmt.Sprintf(":%s", c.GetString("app.port")), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)

}
