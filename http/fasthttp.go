package http

// github.com/ifeisier/go-utils

import (
	"github.com/valyala/fasthttp"
	"time"
)

var (
	client                *fasthttp.Client
	headerContentTypeJson = []byte("application/json;charset=UTF-8")
	headerContentTypeForm = []byte("application/x-www-form-urlencoded")
)

func init() {
	client = createClient()
}

// PostFormRequest 表单(x-www-form-urlencoded) 格式的 post 请求
func PostFormRequest(url string, json []byte, header, cookie map[string]string) (statusCode int32, body []byte, err error) {
	req := acquireRequest(url, fasthttp.MethodPost, json, headerContentTypeForm, header, cookie)
	return sendRequest(req)
}

// PostJsonRequest json 格式的 post 请求
func PostJsonRequest(url string, json []byte, header, cookie map[string]string) (statusCode int32, body []byte, err error) {
	req := acquireRequest(url, fasthttp.MethodPost, json, headerContentTypeJson, header, cookie)
	return sendRequest(req)
}

// GetRequest 一个简单的 Get 请求
func GetRequest(url string) (statusCode int32, body []byte, err error) {
	req := acquireRequest(url, fasthttp.MethodGet, nil, headerContentTypeJson, nil, nil)
	return sendRequest(req)
}

// acquireRequest 获取一个 Request 实例
func acquireRequest(url string, method string, json, contentType []byte, header, cookie map[string]string) *fasthttp.Request {
	// 从池中获取请求
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	for key, value := range header {
		req.Header.Add(key, value)
	}

	for key, value := range cookie {
		req.Header.SetCookie(key, value)
	}

	if contentType != nil {
		req.Header.SetContentTypeBytes(contentType)
	}

	if json != nil {
		req.SetBodyRaw(json)
	}
	return req
}

// sendRequest 真正的发送 http 请求
func sendRequest(req *fasthttp.Request) (statusCode int32, body []byte, err error) {
	// 从池中获取响应
	resp := fasthttp.AcquireResponse()

	// 发出请求，将结果放到 Response 中。
	err = client.Do(req, resp)

	// 用完后释放，减少 GC。
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 状态码和响应的数据
	statusCode = int32(resp.StatusCode())
	body = resp.Body()
	return
}

func createClient() *fasthttp.Client {
	return &fasthttp.Client{
		// 设置 User-Agent
		Name:                     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		NoDefaultUserAgentHeader: true, // 没有指定 User-Agent 就使用 fasthttp
		MaxConnsPerHost:          512,  // 每台主机可建立的最大连接数

		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		WriteBufferSize: 1024 * 1024,
		ReadBufferSize:  1024 * 1024 * 10,

		// 设置 Keep-Alive 的时间
		// 它允许客户端（通常是浏览器）和服务器之间的多个 HTTP 请求和响应共享同一条 TCP 连接，而不是在每个请求之后都建立一个新的连接。
		// 这有助于提高性能，减少延迟，减轻服务器负担，因为不必频繁地建立和关闭连接。
		MaxIdleConnDuration: time.Hour,

		// 禁止 Header 名称规范化
		// 比如 Header 名称 content-type 规范化为 Content-Type
		DisableHeaderNamesNormalizing: true,
		// 禁止 Path 规范化
		// 比如特殊路径 /a//b/../c，规范化为 /a/c
		DisablePathNormalizing: true,

		// 包含了 TCP 连接池，复用 TCP 连接，避免每次请求创建新的 TCP 连接。
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,      // 最大并发数，0表示无限制
			DNSCacheDuration: time.Hour, // DNS 缓存的过期时间
		}).Dial,
	}
}
