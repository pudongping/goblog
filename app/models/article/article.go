package article

import (
	"github.com/pudongping/goblog/app/models"
	"github.com/pudongping/goblog/app/models/user"
	"github.com/pudongping/goblog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body  string `gorm:"type:longtext;not null;" valid:"body"`
	UserID uint64 `gorm:"not null;index"`
	CategoryID uint64 `gorm:"not null;default:0;index"`
	User user.User  // 实现模型关联
}

// Link 方法用来生成文章链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", article.GetStringID())
}

// CreatedAtDate 创建日期
func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}