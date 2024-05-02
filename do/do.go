package do

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type do_allApiFunc func(string)

var doApiMap = map[string]do_allApiFunc{
	"DO_WHOLE": do_get_whole_value,
	"DO_STATUS": do_get_status,
	"DO_PAULSESTATUS": do_get_paulse_status,
	"DO_PAULSECOUNT": do_get_paulse_count,

}

func muti_thread_api(callback func()){
	callback();
}


func Do_choose_api(w http.ResponseWriter, r *http.Request) {
	var apiKey = ""
	fmt.Println("yes this is do")
	w.Write([]byte("Hello, world!"))
	apiKey = "DO_WHOLE"
	if fn, ok := doApiMap[apiKey];ok{
		fn("i'm "+apiKey)
	}else{
		fmt.Println("not exit")
	}
	
	
}
//TODO: get server ip func, finish other api，get slot number


func sendGETRequest(uri string) {
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Accept", "vdn.dac.v1")
	req.Header.Set("Content-Type", "application/json")

	// 使用默认的 HTTP 客户端发送请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 输出响应体
	fmt.Println("Response body:", string(body))
}


func do_get_whole_value(msg string){
	fmt.Println("do_get_whole_value:", msg)
	uri := "http://192.168.127.254/api/slot/0/io/do"
	sendGETRequest(uri)
}

func do_get_status(msg string){
	fmt.Println("do_get_status:", msg)
}

func do_get_paulse_status(msg string){
	fmt.Println("do_get_paulse_status:", msg)
}

func do_get_paulse_count(msg string){
	fmt.Println("do_get_paulse_count:", msg)
}

