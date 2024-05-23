package do

import (
	"fmt"
	http_rest "go-Moxa/http"
	"io/ioutil"
)
var g_nickNameMap = NewNickNameMap()

type NickNameMap struct {
    IndexToValue map[int]string
    ValueToIndex map[string]int
}

func NewNickNameMap() *NickNameMap {
    return &NickNameMap{
        IndexToValue: make(map[int]string),
        ValueToIndex: make(map[string]int),
    }
}

func (n *NickNameMap) AddNickn(value string, index int) {
    n.IndexToValue[index] = value
    n.ValueToIndex[value] = index
}

func (n *NickNameMap) RemoveNickn(value string) {
    if index, ok := n.ValueToIndex[value]; ok {
        delete(n.ValueToIndex, value)
        delete(n.IndexToValue, index)
    }
}

func (n *NickNameMap) GetNicknIndex(value string) (int, bool) {
    index, ok := n.ValueToIndex[value]
    return index, ok
}

func (n *NickNameMap) GetNicknValue(index int) (string, bool) {
    value, ok := n.IndexToValue[index]
    return value, ok
}

func do4510_get_value(apiKey string, machine *Machine, slotName string, chName string) int {
	fmt.Println("do_get_status:", apiKey, " ,ip:", machine.IP, " ,ch:", chName)
	g_nickNameMap.AddNickn(chName, 0)
	do4510_get_rest_request(apiKey, machine, slotName, chName)
	// fmt.Println(res)
	return 0
}



func do4510_get_rest_request(apiKey string, machine *Machine, slotName string, chName string) {

	param := "/doStatus"
	ioName := slotName +"@"+chName
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
