package common

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient 类型
type HTTPClient struct {
	cookieManager map[string]*http.Cookie
	proxy         *url.URL
}

// ClientResponse 表示Http响应
type ClientResponse struct {
	Body    string
	Headers http.Header
}

// NewHTTPClient 表示HTTPClient构造函数
func NewHTTPClient() *HTTPClient {
	c := &HTTPClient{}
	c.cookieManager = make(map[string]*http.Cookie)
	return c
}

func getClient(proxy *url.URL) (client *http.Client) {
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*5) //建立连接超时
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(10 * time.Second)) //发送接收数据超时
			return c, nil
		},
		DisableKeepAlives: true,
	}
	if proxy != nil {
		tr.Proxy = http.ProxyURL(proxy)
	}
	client = &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

// ChangeProxy 用于切换代理
func (c *HTTPClient) ChangeProxy(proxy *url.URL) {
	c.proxy = proxy
}

// HTTPPost 用于发送Post请求
func (c *HTTPClient) HTTPPost(addr string, headers map[string]string, data *url.Values, proxy *url.URL) (result ClientResponse, err error) {
	var (
		bodyReader io.Reader
		req        *http.Request
		resp       *http.Response
	)

	if data != nil {
		bodyReader = strings.NewReader(data.Encode())
	}

	req, err = http.NewRequest("POST", addr, bodyReader)
	req.Close = true

	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.cookieManager {
		req.AddCookie(cookie)
	}

	resp, err = getClient(proxy).Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	c.updateCookie(resp.Cookies())

	result.Headers = resp.Header
	result.Body, err = readResponse(resp)

	return
}

// HTTPStrPost 用于提交string类型的Post请求
func (c *HTTPClient) HTTPStrPost(addr string, headers map[string]string, postData string, proxy *url.URL) (result ClientResponse, err error) {
	var (
		bodyReader io.Reader
		req        *http.Request
		resp       *http.Response
	)

	bodyReader = strings.NewReader(postData)

	req, err = http.NewRequest("POST", addr, bodyReader)
	req.Close = true
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.cookieManager {
		req.AddCookie(cookie)
	}

	resp, err = getClient(proxy).Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	c.updateCookie(resp.Cookies())

	result.Headers = resp.Header
	result.Body, err = readResponse(resp)

	return
}

// HTTPGet 用于Get请求
func (c *HTTPClient) HTTPGet(addr string, headers map[string]string, proxy *url.URL) (result ClientResponse, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	req, err = http.NewRequest("GET", addr, nil)
	req.Close = true
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = getClient(proxy).Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	result.Headers = resp.Header
	result.Body, err = readResponse(resp)
	return
}

// HTTPGetWithParam 用于带参数的Get请求
func (c *HTTPClient) HTTPGetWithParam(addr string, headers map[string]string, query map[string]string, proxy *url.URL) (result ClientResponse, err error) {
	var (
		req  *http.Request
		resp *http.Response
		u    *url.URL
	)
	u, err = url.Parse(addr)
	if err != nil {
		return
	}

	q := u.Query()

	for k, v := range query {
		q.Add(k, v)
	}

	end := strings.Index(addr, "?")

	addr = addr[:end]
	addr = fmt.Sprintf("%s?%s", addr, q.Encode())

	req, err = http.NewRequest("GET", addr, nil)
	req.Close = true
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = getClient(proxy).Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	result.Headers = resp.Header
	result.Body, err = readResponse(resp)
	return
}

// updateCookie 用于更新Cookie
func (c *HTTPClient) updateCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.cookieManager[cookie.Name] = cookie
	}
}

func readResponse(resp *http.Response) (result string, err error) {
	var r io.Reader

	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		r, err = gzip.NewReader(resp.Body)
	} else if strings.Contains(resp.Header.Get("Content-Encoding"), "deflate") {
		r, err = zlib.NewReader(resp.Body)
	} else {
		r = resp.Body
	}

	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	result = string(b)
	return
}
