package do

import (
	"encoding/json"
	"fmt"
	http_rest "go-Moxa/http"
	"io/ioutil"
	"strconv"
)

var do_test_value1 = 0
var do_test_value2 = 0
var do_test_value3 = 0
var do_test_value4 = 0
var do_test_value5 = 0
var do_test_value6 = 0
var do_test_value7 = 0
var do_test_value8 = 0


func NewMachine(main_type string, sub_type string, ip string, ch_type string, numOfChan int) *Machine {
	machine := &Machine{
		Main_type: main_type,
		Sub_type: sub_type,
		IP: ip,
		Ch_type: ch_type,
		Channel:    make([]chan int, numOfChan),
		NumOfChan: numOfChan,
	}
	// 初始化每个通道
	for i := range machine.Channel {
		machine.Channel[i] = make(chan int, 1)
	}
	return machine
}


func (d DoObj)Do_choose_api(apiKey string, machine *Machine, ch int, wrData string) int{
	// var apiKey = ""
	if machine == nil {
        fmt.Println("Machine is nil")
        return -1
    }
	// apiKey = "DO_WHOLE"
	if fn, ok := doApiMap[apiKey];ok{
		return fn(apiKey, machine, ch, wrData)
	}else{
		fmt.Println("not exit")
		return -1
	}

}
//TODO: get server ip func, finish other api，get slot number

func do_update_whole_status(apiKey string, machine *Machine, ch int, wrData string) (map[string]interface{}, error) {
    uri, ok := restUri[apiKey]
    if !ok {
        return nil, fmt.Errorf("API key not found: %s", apiKey)
    }

	if apiKey != "DO_WHOLE" {
		return nil, fmt.Errorf("wrong type")
	}

    fmt.Println("http://" + machine.IP + uri )
    resp, err := http_rest.SendGETRequest("http://" + machine.IP + uri )
    if err != nil {
        return nil, fmt.Errorf("error sending GET request: %s", err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %s", err)
    }
	// fmt.Println("Response body:", string(body))
    var data map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &data); err != nil {
		fmt.Println("Error decoding response body:", err)
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	// for key, value := range data {
	// 	fmt.Println(key, ":", value)
	// }
	doList:= data["io"].(map[string]interface{})["do"].([]interface{})
	for _, item := range doList {
		doItem, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Error: item is not a map[string]interface{}")
			continue
		}
		
		doIndex := doItem["doIndex"].(float64)
		// doMode := doItem["doMode"]
		doStatus := doItem["doStatus"].(float64)
		Do_clear_ch(machine, int(doIndex))
		machine.Channel[int(doIndex)] <- int(doStatus)
	
		// fmt.Println("doIndex:", doIndex, "doMode:", doMode, "doStatus:", doStatus)
	}
    return data, nil
}

func do_get_rest_request(apiKey string, do_param string, machine *Machine, ch int, wrData string)(float64){
	uri, ok := restUri[apiKey]
    if !ok {
        return -1
    }
	param, ok := restParam[do_param]
    if !ok {
        fmt.Println("Error: msg not found in restUri")
        return -1
    }
	fmt.Println("http://"+machine.IP+uri+"/"+strconv.Itoa(ch)+param)
	resp, err := http_rest.SendGETRequest("http://"+machine.IP+uri+"/"+strconv.Itoa(ch)+param)
	if err != nil {
        fmt.Println("Error sending GET request:", err)
        return -1
    }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return -1
	}
	
	// fmt.Println("Response body:", string(body))
	var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        return -1
    }
	io := data["io"].(map[string]interface{})
    do := io[machine.Ch_type].(map[string]interface{})
    status := do[strconv.Itoa(ch)].(map[string]interface{})["doStatus"].(float64)

    return status
}


func do_put_rest_request(apiKey string, do_param string, machine *Machine, ch int, wrData string)(float64){
	uri, ok := restUri[apiKey]
	reqJson := `{"slot":0,"io":{"do":{"`+strconv.Itoa(ch)+`":{"doStatus":`+wrData+`}}}}`
	fmt.Println(reqJson)
    if !ok {
		Do_clear_ch(machine, ch)
        return -1
    }
	param, ok := restParam[do_param]
    if !ok {
        fmt.Println("Error: msg not found in restUri")
		Do_clear_ch(machine, ch)
        return -1
    }
	fmt.Println("http://"+machine.IP+uri+"/"+strconv.Itoa(ch)+param)
	_, err := http_rest.SendPUTRequest("http://"+machine.IP+uri+"/"+strconv.Itoa(ch)+param, reqJson)
	if err != nil {
        fmt.Println("Error sending GET request:", err)
		Do_clear_ch(machine, ch)
        return -1
    }
	num, err := strconv.Atoi(wrData)
	if err != nil {
		fmt.Println("轉換失敗:", err)
		Do_clear_ch(machine, ch)
		return -1
	}

	Do_push_ch(machine, ch, num)

	fmt.Println("success")
	return 0

	
}

func do_get_whole_value(msg string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 0)
	// uri := "http://192.168.127.254/api/slot/0/io/do"
	// http_rest.SendGETRequest(uri)
	return 0
}


func do_get_status(msg string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 1)
	return 0
}

func do_get_paulse_status(msg string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 2)
	return 0
}

func do_get_paulse_count(msg string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 3)
	return 0
}

func _do_check_value(apiKey string, machine *Machine, ch int, wrData string) int{

	value := <-machine.Channel[ch]
	fmt.Println("ch ",ch, ":", value)
	machine.Channel[ch] <- value
	i, err := strconv.Atoi(wrData)
    if err != nil {
        fmt.Println("trans fail:", err)
        return -1
    }
	if i == value{
		return 1
	}else{
		return -1
	}

}



//commom func

func do_check_value(apiKey string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", ch)
	
	return _do_check_value(apiKey, machine, ch, wrData)
}
func do_update_value(apiKey string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", ch)
	do_update_whole_status("DO_WHOLE", machine, 0, wrData)
	
	return 0
}
func do_get_value(apiKey string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", ch)
	res :=do_get_rest_request(apiKey,"DO_STATUS", machine, ch, wrData)
	fmt.Println(res)
	return 0
}

func do_put_value(apiKey string, machine *Machine, ch int, wrData string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", ch)
	// Do_pop_ch(machine, ch)
	Do_clear_ch(machine, ch)
	do_put_rest_request(apiKey,"DO_STATUS", machine, ch, wrData)
	return 0
}

func Do_push_ch(machine *Machine, ch int, value int) int {
	
	fmt.Println("PUSH ip:", machine.IP, " ,ch:", ch, ", value", value)
	// fmt.Println("[####result####]Get value:", )
	select {
		case machine.Channel[ch] <- value:
			return value
		default:
			return -1 
	}

}

func Do_pop_ch(machine *Machine, ch int) int {
	
	
	// fmt.Println("[####result####]Get value:", )
	select {
		case res := <-machine.Channel[ch]:
			fmt.Println("POP ip:", machine.IP, " ,ch:", ch, ", value", res)
			return res
		default:

			return -1 
	}
}

func Do_clear_ch(machine *Machine, ch int) int{
	select {
    case <-machine.Channel[ch]:
        // fmt.Println("ch ", ch, " is non empty, clear it")
		return 0
    default:
        // fmt.Println("ch", ch, " is empty")
		return 1
    }
}


