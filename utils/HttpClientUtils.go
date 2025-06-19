package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

type HttpClientUtils struct {
	Header map[string][]string
}

func (svc *HttpClientUtils) Get(url string, timeout int) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return svc.createRequestError(err)
	}

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PostParams(url string, params string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PostJson(url string, body string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return svc.httpRequest(req, timeout)
}

// 如果需要更通用的方法，可以这样写：
func (svc *HttpClientUtils) PostJsonWithHeaders(url string, body string, headers map[string]string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}

	// 设置基本 header
	req.Header.Set("Content-type", "application/json")

	// 设置额外的 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PutJson(url string, body string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}

	req.Header.Set("Content-type", "application/json")
	svc.setRequestHeader(req)

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PostParamsWithToken(url string, params string, token string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PutParamsWithToken(url string, params string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) PostJsonWithToken(url string, body string, timeout int) ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return svc.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return svc.httpRequest(req, timeout)
}

func (svc *HttpClientUtils) httpRequest(req *http.Request, timeout int) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	svc.setRequestHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header

	return wrapper
}

func (svc *HttpClientUtils) setRequestHeader(req *http.Request) {
	for key, value := range svc.Header {
		req.Header[key] = value
	}
}

func (svc *HttpClientUtils) createRequestError(err error) ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return ResponseWrapper{0, errorMessage, make(http.Header)}
}

// 检测是否返回200
func (svc *HttpClientUtils) CheckUrlHttpCodeOk(url, method string) (bool, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		return false, err
	}
	resp := svc.httpRequest(req, 30)
	ok := resp.StatusCode == http.StatusOK
	if !ok {
		//fmt.Println(resp.Body)
	}
	return ok, errors.New(fmt.Sprintf("http状态码：%d", resp.StatusCode))
}

func (svc *HttpClientUtils) GetValFromQuery(queryStr, key string) (string, error) {
	u, err := url.Parse(queryStr)
	if err != nil {
		return "", errors.New("参数不合法。")
	}

	m, _ := url.ParseQuery(u.RawQuery)

	if len(m[key]) > 0 {
		return m[key][0], nil
	}
	return "", nil
}

func (svc *HttpClientUtils) DownloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// 创建本地文件
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// 将响应写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
