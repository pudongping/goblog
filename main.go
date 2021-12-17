package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/pudongping/goblog/app/http/middlewares"
	"github.com/pudongping/goblog/bootstrap"
	"github.com/pudongping/goblog/pkg/database"
	"github.com/pudongping/goblog/pkg/logger"
)

var router *mux.Router
var db *sql.DB

func main() {

	database.Initialize()
	db = database.DB

	// 初始化 gorm
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)

	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
