package requests

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thedevsaddam/govalidator"

	"github.com/pudongping/goblog/pkg/model"
)

// 此方法会在 requests 包引入时执行
func init()  {
	// 这里自定义了验证规则
	// not_exists:users,email
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbFiled+" = ?", val).Count(&count)

		if count != 0 {

			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}

		return nil
	})
}