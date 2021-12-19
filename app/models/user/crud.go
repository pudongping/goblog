package user

import (
	"github.com/pudongping/goblog/pkg/logger"
	"github.com/pudongping/goblog/pkg/model"
	"github.com/pudongping/goblog/pkg/types"
)

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// Get 通过 ID 获取用户
func Get(idstr string) (User, error) {
	var user User
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil
}

// GetByEmail 通过 email 来获取用户
func GetByEmail(email string) (User, error) {
	var user User
	err := model.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
