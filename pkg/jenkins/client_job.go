package jenkins

import (
	"reflect"
	"regexp"
)

// JobClient 封裝了跟 job 有關的 jenkins rest api
type JobClient struct {
	c *Client
}

// Job 取得 JobClient
func (c *Client) Job() *JobClient {
	return &JobClient{
		c: c,
	}
}

// Job 代表一個 jenkins 上的 job 資料
type Job struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
}

// Match 判斷傳入的 filter 是否用 regex 符合此 struct 任一欄位
func (j Job) Match(filter string) bool {
	if len(filter) == 0 {
		return true
	}
	expr, err := regexp.Compile(filter)
	if err != nil {
		return false
	}
	v := reflect.ValueOf(j)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if s, ok := field.Interface().(string); ok && expr.MatchString(s) {
			return true
		}
	}
	return false
}
