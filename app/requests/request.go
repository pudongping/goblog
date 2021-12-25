package requests

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/thedevsaddam/govalidator"

	"github.com/pudongping/goblog/pkg/model"
)

// 此方法会在 requests 包引入时执行
func init() {
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

	// 因为中文文字的长度判断需要使用 utf8.RuneCountInString() 来验证，因此这里需要重新定义长度判断

	// max_cn:8
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}

		return nil
	})

	// min_cn:2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}

		return nil
	})

}
