package view

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/pudongping/goblog/pkg/logger"
	"github.com/pudongping/goblog/pkg/route"
)

// Render 渲染视图
func Render(w io.Writer, data interface{}, tplFiles ...string) {
	// 1 设置模板相对路径
	viewDir := "resources/views/"

	// 2. 语法糖，将 articles.show 更正为 articles/show
	// 第四个参数表示允许替换的次数，设置为 -1 表示替换所有
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	// 3 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 4 在 Slice 里新增我们的目标文件
	allFiles := append(layoutFiles, tplFiles...)

	// 5 解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 6 渲染模板
	err = tmpl.ExecuteTemplate(w, "myapp", data)
	logger.LogError(err)
}
