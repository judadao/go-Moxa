package do

import (
	"encoding/json"
	"fmt"
	http_rest "go-Moxa/http"
	"io/ioutil"
	"strconv"
)
var g_nickChMap = NewNickChMap()

type ChNickMap struct {
    NickToValue map[string]int
    NickToIndex map[string]int
	currentIndex int
}

func NewNickChMap() *ChNickMap {
    return &ChNickMap{
        NickToValue: make(map[string]int),
        NickToIndex: make(map[string]int),
		currentIndex: 0,
    }
}

func (n *ChNickMap) AddNickn(nickname string) {
    n.NickToValue[nickname] = 0
    n.NickToIndex[nickname] = n.currentIndex
	n.currentIndex++
}

func (n *ChNickMap) RemoveNickn(nickname string) {

    delete(n.NickToValue, nickname)
    delete(n.NickToIndex, nickname)
}

func (n *ChNickMap) GetNicknIndex(nickName string) (int, bool) {
    value, ok := n.NickToIndex[nickName]
    return value, ok
}

func (n *ChNickMap) GetNicknValue(nickName string) (int, bool) {
    value, ok := n.NickToValue[nickName]
    return value, ok
}

func do4510_check_nick_map(chName string) int {
	if _, exists := g_nickChMap.NickToValue[chName]; !exists {
        g_nickChMap.AddNickn(chName)
        // fmt.Println("Key added:", chName)
    } else {
        return 0
    }
	return 0

}

func do4510_get_value(apiKey string, machine *Machine, chName string, value string) int {
	// fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", chName, ",value:", value)
	do4510_check_nick_map(chName)
	// num, err := strconv.Atoi(value)
	// if err != nil {
	// 	fmt.Println("轉換失敗:", err)
	// 	Do_clear_ch(machine, g_nickChMap.NickToIndex[chName])
	// 	return -1
	// }
	num := do4510_get_rest_request(apiKey, machine, chName)
	if do4510_get_rest_request(apiKey, machine, chName) == -1{
		return -1
	}
	if index, exists := g_nickChMap.NickToIndex[chName]; exists {
		
		Do_clear_ch(machine, index)
		Do_push_ch(machine, index, num)
	} else {
		fmt.Println("Error: Index for", chName, "not found in NickToIndex map")
		return -1
	}
	
	// fmt.Println(res)
	return 0
}

func do4510_put_value(apiKey string, machine *Machine, chName string, value string) int {
	// fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", chName)
	do4510_check_nick_map(chName)
	
	num, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("轉換失敗:", err)
		Do_clear_ch(machine, g_nickChMap.NickToIndex[chName])
		return -1
	}
	if do4510_put_rest_request(apiKey, machine, chName, value) == -1{
		return -1
	}
	
	if index, exists := g_nickChMap.NickToIndex[chName]; exists {
		Do_clear_ch(machine, index)
		Do_push_ch(machine, index, num)
	} else {
		fmt.Println("Error: Index for", chName, "not found in NickToIndex map")
		return -1
	}

	return 0
}



func do4510_get_rest_request(apiKey string, machine *Machine, chName string) int {

	param := "/doStatus"
	ioName := machine.Slot_nick +"@"+chName
	uri := "http://" + machine.IP + "/api" + "/io/do/" + ioName + param
	// fmt.Println(uri)
	resp, err := http_rest.SendGETRequest_v2(uri)
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
	// io := data["io"].(map[string]interface{})
	// do := io[machine.Ch_type].(map[string]interface{})
	status, ok := data["value"].(float64)
	if !ok {
		fmt.Println("Error: 'value' is not a float64")
		return -1
	}
	
	statusInt := int(status)
	// fmt.Println("status", statusInt)
	g_nickChMap.NickToValue[chName] = statusInt
	// fmt.Println("status", chName)


	return statusInt
}


func do4510_put_rest_request(apiKey string, machine *Machine, chName string, value string) int {

	param := "/doStatus"
	ioName := machine.Slot_nick +"@"+chName
	uri := "http://" + machine.IP + "/api" + "/io/do/" + ioName + param
	// fmt.Println(uri)
	reqJson := `{"value":`+value+`}`
	resp, err := http_rest.SendPUTRequest_v2(uri, reqJson)
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
	// io := data["io"].(map[string]interface{})
	// do := io[machine.Ch_type].(map[string]interface{})
	status, ok := data["value"].(float64)
	if !ok {
		fmt.Println("Error: 'value' is not a float64")
		return -1
	}
	// fmt.Println("status", status)
	statusInt := int(status)
	


	return statusInt
}



func _do4510_check_value(apiKey string, machine *Machine, chName string, wrData string) int {

	var ch, _=g_nickChMap.GetNicknIndex(chName)
	value := <-machine.Channel[ch]
	// fmt.Println("ch ", ch, ":", value)
	machine.Channel[ch] <- value
	i, err := strconv.Atoi(wrData)
	if err != nil {
		fmt.Println("trans fail:", err)
		return -1
	}
	// fmt.Println("_do4510_check_value wrData:", wrData)
	// fmt.Println("_do4510_check_value value:", value)
	if i == value {
		return 1
	} else {
		return -1
	}

}

func do4510_check_value(apiKey string, machine *Machine,  chName string, wrData string) int {
	// fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", chName)

	return _do4510_check_value(apiKey, machine, chName, wrData)
}