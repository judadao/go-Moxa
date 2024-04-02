package main

import (
	"fmt"
	"io"
	"net/http"
)

type DiHandler struct {
}

func resErr(ip string, w io.Writer, headers map[string]string) error {
    // 創建新的 HTTP 請求
	url := "http://" + ip + "/api/slot/0/io/di"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    // 設置自定義標頭
    for key, value := range headers {
        req.Header.Set(key, value)
    }

    // 發送請求
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 將回應寫入到指定的 io.Writer
    _, err = fmt.Fprintln(w, "HTTP 請求成功，回應內容為:")
    if err != nil {
        return err
    }

    _, err = io.Copy(w, resp.Body)
    if err != nil {
        return err
    }

    return nil
}


func (d DiHandler) Get_di(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
    if ip == "" {
        http.Error(w, "Missing IP Address", http.StatusBadRequest)
        return
    }
	headers := map[string]string{
        "Accept": "vdn.dac.v1",
        "Content-Type":  "application/json",
        // 添加其他所需的標頭
    }
        err := resErr(ip, w,headers)
        if err != nil {
            fmt.Println("HTTP 請求錯誤:", err)
            return
        }
}

