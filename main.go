package main

import (
	"embed"
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

// embed 读取相对路径是相对于书写 //go:embed 指令的 .go 文件的
// go:embed 我们称之为指令，后面跟着文件目录，或者是单个文件，支持通配符。public/* 意味着递归加载 public 下所有子目录以及文件
// embed.FS 就是一个打包在二进制包里的文件系统；
// //go:embed 指令是告知要加载哪些文件目录；
// 加载和读取目录时，会 保持原有的目录结构；
// 除了写功能，你可以将其当做一般文件系统来读取里面的文件和目录。

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {

	//database.Initialize()
	//db = database.DB

	// 初始化 gorm
	bootstrap.SetupDB()

	// 初始化模版
	bootstrap.SetupTemplate(tplFS)

	// 初始化路由绑定
	router = bootstrap.SetupRoute()

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)

	err := http.ListenAndServe(fmt.Sprintf(":%s", c.GetString("app.port")), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)

}
