package util

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

var CopilotHeaders = map[string]string{
	"Accept":                "*/*",
	"Accept-Encoding":       "gzip,deflate,br",
	"Connection":            "close",
	"Editor-Plugin-Version": "copilot-intellij/1.4.4.3955",
	"Editor-Version":        "JetBrains-WS/233.11799.229",
	"User-Agent":            "GithubCopilot/1.140.0",
}

func SendHTTPRequest(method string, url string, headers map[string]string, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b := HandleResponseEncoding(resp)

	return b, nil
}

func HandleResponseEncoding(resp *http.Response) []byte {
	contentEncoding := resp.Header.Get("Content-Encoding")
	var reader io.Reader
	switch contentEncoding {
	case "br":
		brReader := brotli.NewReader(resp.Body)
		reader = brReader
	case "gzip":
		gzReader, _ := gzip.NewReader(resp.Body)
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		flateReader := flate.NewReader(resp.Body)
		defer flateReader.Close()
		reader = flateReader
	default:
		reader = resp.Body
	}
	responseData, _ := io.ReadAll(reader)
	return responseData
}

func SendBarkNotice(title, content string) {
	_, _ = SendHTTPRequest("GET", "http://182.43.242.4:7890/zxuNx7uD6Q2AtoJgSVmFqc"+"/"+title+"/"+content+"?icon=http://182.43.242.4:3000/dist/favicon.png", nil, nil)
}