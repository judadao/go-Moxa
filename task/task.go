package allTask

import (
	"fmt"
	"go-Moxa/do"
	task_qu "go-Moxa/queue"
	"strconv"
)

func Task_func(taskInfo *task_qu.Task) {
	// var lastCheckResult int
	doObj := do.DoObj{}
	do := do.NewMachine(taskInfo.Device.Type, taskInfo.Device.SubType, taskInfo.Device.IP, "do", 8)
	doObj.Do_choose_api("DO_WHOLE", do, 0, "")
	res := 1
	for _, condi := range taskInfo.When { // 遞迴比較條件，皆為1才會過
		chnum, _ := strconv.Atoi(condi.ChNum)
		checkResult := doObj.Do_choose_api("DO_CHECK", do, chnum, condi.Value)
		// fmt.Println("[##res##]", checkResult)
		if checkResult == -1 {
			res &= 0
		} else {
			res &= 1
		}

	}

	if res == taskInfo.LastCheck {
		// check last
		fmt.Println("Status is the same as last check, no action needed.")
		return
	}
	taskInfo.LastCheck = res // update status
	if res == 1 {
		for _, action := range taskInfo.ThenActions {
			chnum, _ := strconv.Atoi(action.ChNum)
			doObj.Do_choose_api("DO_PUT_VALUE", do, chnum, action.Value)
		}
	}

}