package common

import (
	"testing"
)

func Test12306InitGet(t *testing.T) {
	c := NewKyfwClient()
	headers := make(map[string]string)
	headers["Cache-Control"] = "max-age=0"
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	headers["Upgrade-Insecure-Requests"] = "1"
	resp, err := c.Get("https://kyfw.12306.cn/otn/login/init", headers)
	if err != nil {
		t.Error("请求失败:", err)
	} else {
		t.Log("请求成功.")
		t.Log("request header:", c.header)
		t.Log("response header:", resp.header)
		t.Log("response body:", resp.body)
	}
}
