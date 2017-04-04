package common

import (
	"compress/gzip"

	"io/ioutil"
	"net/http"
)

type KyfwClient struct {
	header map[string]string
	client *http.Client
}

type KyfwResponse struct {
	body   string
	header http.Header
}

func NewKyfwClient() *KyfwClient {
	baseHeader := make(map[string]string)
	baseHeader["Accept-Encoding"] = "gzip, deflate, sdch, br"
	baseHeader["Accept-Language"] = "zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4"
	baseHeader["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
	baseHeader["AlexaToolbar-ALX_NS_PH"] = "AlexaToolbar/alx-4.0.1"
	baseHeader["Connection"] = "keep-alive"
	baseHeader["Host"] = "kyfw.12306.cn"

	return &KyfwClient{header: baseHeader, client: &http.Client{}}
}

func (c *KyfwClient) Get(url string, headers map[string]string) (*KyfwResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		if _, ok := c.header[k]; !ok {
			c.header[k] = v
		}
	}

	for k, v := range c.header {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	r := &KyfwResponse{}
	r.header = resp.Header
	var bt []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		bt, err = ioutil.ReadAll(reader)
	default:
		bt, err = ioutil.ReadAll(resp.Body)
	}

	r.body = string(bt)
	return r, err
}

func (c *KyfwClient) Post(url string, headers map[string]string) (*KyfwResponse, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		if _, ok := c.header[k]; !ok {
			c.header[k] = v
		}
	}

	for k, v := range c.header {
		req.Header.Add(k, v)
	}
	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	r := &KyfwResponse{}
	r.header = resp.Header
	var bt []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		bt, err = ioutil.ReadAll(reader)
	default:
		bt, err = ioutil.ReadAll(resp.Body)
	}

	r.body = string(bt)
	return r, err
}
