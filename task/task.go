package allTask

import (
	"fmt"
	"go-Moxa/do"
	task_qu "go-Moxa/queue"
	"time"
)

// TODO: add 4510 判斷device type，來分別用choose api，因為都已轉換成string所以可以直接判斷 type得取即可，obj 可能比較麻煩
func Task_func(taskInfo *task_qu.Task) {
	// var lastCheckResult int
	
	doObjs := make(map[string]*do.Machine)  // creat obj map
	// fmt.Println(taskInfo.Device)

	for name, device := range taskInfo.Device {
		// fmt.Println("[####]Number of conditions:", name)
		//4510 need to use subtype to set obj
		if device.Type == "1200" {
			doObjs[name] = do.NewMachine(device.Type, device.SubType, device.IP, "do", do.Subtype_map[device.SubType])
		}else {
			doObjs[name] = do.NewMachine_4510(device.Type, device.SubType, device.IP, device.Name, "do", do.Subtype_map[device.SubType])
		}
		
	}

	res := 1
	for _, condi := range taskInfo.When {
		
		deviceName := condi.ObjName

		
		if doObj, exists := doObjs[deviceName]; exists {
			chnum := condi.ChNum
			chval := condi.Value
			do.Do_choose_api("DO_GET_VALUE", doObj, chnum, chval)
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
		// fmt.Println("Status is the same as last check, no action needed.")
		return
	}
	retryCount := 5
	taskInfo.LastCheck = res // update status
	if res == 1 {
        for _, action := range taskInfo.ThenActions {
            deviceName := action.ObjName
			// fmt.Println("Processing condition %d: Device = %s\n", action)
            if doObj, exists := doObjs[deviceName]; exists {
                chnum:= action.ChNum
				for i := 0; i < retryCount; i++ {

					result := do.Do_choose_api("DO_PUT_VALUE", doObj, chnum, action.Value)

					if result != -1 {
						// fmt.Println("Do_choose_api successful!")
						break
					}

					if i == retryCount-1 {
						fmt.Println("Do_choose_api failed after", retryCount, "attempts")
					}

					time.Sleep(1 * time.Second)
				}
            } else {
                fmt.Printf("Device %s not found.\n", deviceName)
            }
        }
    }

}



