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


func (d DoObj)Do_choose_api(apiKey string, machine *Machine, ch int) int{
	// var apiKey = ""
	if machine == nil {
        fmt.Println("Machine is nil")
        return -1
    }
	// apiKey = "DO_WHOLE"
	if fn, ok := doApiMap[apiKey];ok{
		return fn(apiKey, machine, ch)
	}else{
		fmt.Println("not exit")
		return -1
	}

}
//TODO: get server ip func, finish other api，get slot number

func do_get_rest_request(apiKey string, do_param string, machine *Machine, ch int)(float64){
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

func do_get_whole_value(msg string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 0)
	// uri := "http://192.168.127.254/api/slot/0/io/do"
	// http_rest.SendGETRequest(uri)
	return 0
}


func do_get_status(msg string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 1)
	return 0
}

func do_get_paulse_status(msg string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 2)
	return 0
}

func do_get_paulse_count(msg string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_push_ch(machine, ch, 3)
	return 0
}




//commom func
func do_get_value(apiKey string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", ch)
	res :=do_get_rest_request(apiKey,"DO_STATUS", machine, ch)
	fmt.Println(res)
	return 0
}

func do_put_value(msg string, machine *Machine, ch int) int {
	fmt.Println("do_get_status:", msg, " ,ip:", machine.IP, " ,ch:", ch)
	Do_pop_ch(machine, ch)
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
        fmt.Println("ch is non empty, clear it")
		return 0
    default:
        fmt.Println("ch is empty")
		return 1
    }
}


