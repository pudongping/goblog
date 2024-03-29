package controllers

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/pudongping/goblog/app/models/article"
	"github.com/pudongping/goblog/app/policies"
	"github.com/pudongping/goblog/app/requests"
	"github.com/pudongping/goblog/pkg/auth"
	"github.com/pudongping/goblog/pkg/logger"
	"github.com/pudongping/goblog/pkg/route"
	"github.com/pudongping/goblog/pkg/view"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
	BaseController
}

// Index 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	// 1. 获取结果集
	// articles, err := article.GetAll()
	articles, pagerData, err := article.GetAll(r, 2)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 2. 加载模版
		view.Render(w, view.D{
			"Articles": articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}

}

// Show 文章详情页面
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 4. 读取成功，显示文章
		view.Render(w, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article), // 权限判断按钮
		}, "articles.show", "articles._article_meta")
	}

}

// Create 文章创建页面
func (ac *ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
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
func (ac *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	currentUser := auth.User()

	// 1. 初始化数据
	_article := article.Article{
		Title:  title,
		Body:   body,
		UserID: currentUser.ID,
	}

	// 2. 表单验证
	errors := requests.ValidateArticleForm(_article)

	// 检查是否含有错误
	if len(errors) == 0 {
		// 创建文章
		_article.Create()

		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
			// 第二个参数表示为转换为十进制
			// fmt.Fprint(w, "插入成功， ID 为"+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "文章创建失败，500 服务器内部错误")
		}
	} else {

		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")

	}

}

// Edit 文章更新页面
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 4. 读取成功，显示编辑文章表单
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  view.D{},
			}, "articles.edit", "articles._form_field")
		}

	}

}

// Update 更新文章
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 4. 未出现错误

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {

			// 4.1 表单验证
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			errors := requests.ValidateArticleForm(_article)

			if len(errors) == 0 {

				// 4.2 表单验证通过，更新数据
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
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  errors,
				}, "articles.edit", "articles._form_field")

			}

		}

	}

}

// Delete 删除文章
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
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

}
