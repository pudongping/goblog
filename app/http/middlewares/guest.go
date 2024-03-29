package middlewares

import (
	"net/http"

	"github.com/pudongping/goblog/pkg/auth"
	"github.com/pudongping/goblog/pkg/flash"
)

// Guest 只允许未登录用户访问
func Guest(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if auth.Check() {
			flash.Warning("登录用户无法访问此页面")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// 继续处理接下去的请求
		next(w, r)
	}
}
