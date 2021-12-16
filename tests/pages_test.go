package tests

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPages(t *testing.T) {

	baseURL := "http://localhost:3000"

	// 1. 声明加载初始化测试数据
	var tests = []struct {
		method   string  // 请求方法
		url      string  // URI
		expected int  // 状态码
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/100/delete", 404},
	}

	// 2. 遍历所有测试
	for _, test := range tests {
		t.Logf("当前请求 URL：%v \n", test.url)
		var (
			resp *http.Response
			err  error
		)

		// 2.1 请求以获取响应
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}

		// 2.2 断言
		assert.NoError(t, err, fmt.Sprintf("请求 %s 时报错", test.url))
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" 应该返回状态码"+strconv.Itoa(test.expected))

	}

}
