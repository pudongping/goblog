/**
文章的权限判断
*/
package policies

import (
	"github.com/pudongping/goblog/app/models/article"
	"github.com/pudongping/goblog/pkg/auth"
)

// CanModifyArticle 是否允许修改话题
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
