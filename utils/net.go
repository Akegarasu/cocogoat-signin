package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2:   true,
			MaxConnsPerHost:     0,
			MaxIdleConns:        0,
			MaxIdleConnsPerHost: 999,
		},
	}

	// UserAgent HTTP请求时使用的UA
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66"
)

func PostBytes(url string, data []byte, headers map[string]string) ([]byte, error) {
	reader, err := HTTPPostReadCloser(url, data, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}

func GetBytes(url string, headers map[string]string) ([]byte, error) {
	reader, err := HTTPGetReadCloser(url, headers)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return ioutil.ReadAll(reader)
}

// HTTPGetReadCloser 从 Http url 获取 io.ReadCloser
func HTTPGetReadCloser(url string, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["User-Agent"]; !ok {
		req.Header["User-Agent"] = []string{UserAgent}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

// HTTPPostReadCloser 从 Http url 获取 io.ReadCloser
func HTTPPostReadCloser(url string, data []byte, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["User-Agent"]; !ok {
		req.Header["User-Agent"] = []string{UserAgent}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
