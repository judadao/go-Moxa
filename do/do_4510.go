package do

import (
	"fmt"
	http_rest "go-Moxa/http"
	"io/ioutil"
	"strconv"
)
var g_nickChMap = NewNickChMap()

type ChNickMap struct {
    NickToValue map[string]int
    NickToIndex map[string]int
}

func NewNickChMap() *ChNickMap {
    return &ChNickMap{
        NickToValue: make(map[string]int),
        NickToIndex: make(map[string]int),
    }
}

func (n *ChNickMap) AddNickn(nickname string, index int, value int) {
    n.NickToValue[nickname] = value
    n.NickToIndex[nickname] = index
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

func do4510_get_value(apiKey string, machine *Machine, chName string, value string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", chName)
	g_nickChMap.AddNickn(chName, 0, 0)
	num, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("轉換失敗:", err)
		Do_clear_ch(machine, g_nickChMap.NickToIndex[chName])
		return -1
	}
	Do_push_ch(machine, g_nickChMap.NickToIndex[chName], num)
	do4510_get_rest_request(apiKey, machine, chName)
	// fmt.Println(res)
	return 0
}



func do4510_get_rest_request(apiKey string, machine *Machine, chName string) {

	param := "/doStatus"
	ioName := machine.Slot_nick +"@"+chName
	uri := "http://" + machine.IP + "/api" + "/io/do/" + ioName + param
	fmt.Println(uri)
	resp, err := http_rest.SendGETRequest_v2(uri)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		// return -1
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		// return -1
		return
	}

	fmt.Println("Response body:", string(body))
	// var data map[string]interface{}
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	return -1
	// }
	// io := data["io"].(map[string]interface{})
	// do := io[machine.Ch_type].(map[string]interface{})
	// status := do[strconv.Itoa(ch)].(map[string]interface{})["doStatus"].(float64)

	// return status
}
