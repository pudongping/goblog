package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"

	"gorm.io/gorm"

	"github.com/pudongping/goblog/app/models/article"
	"github.com/pudongping/goblog/pkg/logger"
	"github.com/pudongping/goblog/pkg/route"
	"github.com/pudongping/goblog/pkg/types"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	// 1. 获取结果集
	articles, err := article.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 2. 加载模版
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		// 3. 渲染模版，将所有文章的数据传输进去
		// Execute() 在执行时会设置正确的 HTML 标头
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}

}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}

}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	// 可以使用类似于以下的语法修改默认的模版标识符，比如这里将默认的 {{}} 修改成 {[]}
	// template.New("test").Delims("{[", "]}").ParseFiles("filename.gohtml")
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)

	titleLen := utf8.RuneCountInString(title) // 计算 title 的长度
	bodyLen := utf8.RuneCountInString(body)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if titleLen < 3 || titleLen > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if bodyLen < 10 {
		errors["body"] = "内容长度需要大于或等于 10 个字节"
	}

	return errors
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否含有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		// 创建文章
		_article.Create()

		if _article.ID > 0 {
			// 第二个参数表示为转换为十进制
			fmt.Fprint(w, "插入成功， ID 为"+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		storeURL := route.Name2URL("articles.store")

		// 构建 ArticlesFormData 里的数据，storeURL 是通过路由参数生成的 URL 路径
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)

	}

}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示编辑文章表单
		updateURL := route.Name2URL("articles.update", "id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}

}

// Update 更新文章
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误

		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {

			// 4.2 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()

			if err != nil {
				// 数据库错误
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			// 更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}

		} else {

			// 4.3 表单验证不通过，显示理由

			updateURL := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			err = tmpl.Execute(w, data)
			logger.LogError(err)

		}

	}

}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request)  {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := _article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 sql 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}

	}

}
