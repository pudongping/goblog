package category

import (
	"github.com/pudongping/goblog/pkg/logger"
	"github.com/pudongping/goblog/pkg/model"
)

// Create 创建分类，通过 category.ID 来判断是否创建成功
func (category *Category) Create() (err error) {
	if err = model.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}
