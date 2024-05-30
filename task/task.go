package allTask

import (
	"fmt"
	"go-Moxa/do"
	task_qu "go-Moxa/queue"
)

func Task_func(taskInfo *task_qu.Task) {
	// var lastCheckResult int
	
	doObjs := make(map[string]*do.Machine)  // creat obj map
	for name, device := range taskInfo.Device {
		// fmt.Println("[####]Number of conditions:", name)
		doObjs[name] = do.NewMachine(device.Type, device.SubType, device.IP, "do", do.Subtype_map[device.SubType])
	}

	res := 1
	for _, condi := range taskInfo.When {
		
		deviceName := condi.ObjName

		
		if doObj, exists := doObjs[deviceName]; exists {
			chnum := condi.ChNum
			do.Do_choose_api("DO_GET_VALUE", doObj, chnum, "")
			checkResult := do.Do_choose_api("DO_CHECK", doObj, chnum, condi.Value)
			if checkResult == -1 {
				res &= 0
			} else {
				res &= 1
			}
		} else {
			fmt.Printf("Device %s not found.\n", deviceName)
			res &= 0
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
            deviceName := action.ObjName
			// fmt.Printf("Processing condition %d: Device = %s\n", i, deviceName)
            if doObj, exists := doObjs[deviceName]; exists {
                chnum:= action.ChNum
                do.Do_choose_api("DO_PUT_VALUE", doObj, chnum, action.Value)
            } else {
                fmt.Printf("Device %s not found.\n", deviceName)
            }
        }
    }

}