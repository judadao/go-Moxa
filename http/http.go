package http_rest

import (
	"fmt"
	"net/http"
	"strings"
)

func SendGETRequest(uri string) (*http.Response, error){
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Accept", "vdn.dac.v1")
	req.Header.Set("Content-Type", "application/json")

	// 使用默认的 HTTP 客户端发送请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	return resp, nil
}

func SendPUTRequest(uri string, jsonBody string) (*http.Response, error) {
    // 创建一个新的 HTTP 请求
    req, err := http.NewRequest("PUT", uri, strings.NewReader(jsonBody))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return nil, err
    }

    // 设置请求头
    req.Header.Set("Accept", "vdn.dac.v1")
    req.Header.Set("Content-Type", "application/json")

    // 使用默认的 HTTP 客户端发送请求
    client := http.DefaultClient
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return nil, err
    }

    return resp, nil
}