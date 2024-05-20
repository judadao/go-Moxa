package main

import (
	"fmt"
	"go-Moxa/do"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// var globalDoObjects []*do.Machine

type TaskQueue struct {
	queue []func(task Task)
	mutex sync.Mutex
}

type Task struct {
	Device      DeviceInfo
	When        []Condition
	ThenActions []Action
	LastCheck   int //check last status
}

type DeviceInfo struct {
	Type string
	IP   string
}

type Condition struct {
	ChType string
	ChNum  string
	Value  string
}

type Action struct {
	ChType string
	ChNum  string
	Value  string
}

var wg sync.WaitGroup

func main() {
	taskQueue := NewTaskQueue()

	// var lastCheckResult int
	// taskQueue := make(chan Task, 10)

	// doObj := do.DoObj{}
	// do := do.NewMachine("e1200", "1213", "192.168.127.254", "do", 8)

	// task2(doObj, do)

	// di2 := do.NewMachine("e1200", "1213","192.1.1.2", "di", 8)
	// di3 := do.NewMachine("e1200", "1213","192.1.1.3", "di", 8)
	// di4 := do.NewMachine("e1200", "1213","192.1.1.4", "di", 8)
	// router.Get("/", doObj.Do_choose_api)

	// for i := 0; i < 4; i++ {
	// 	wg.Add(1) // 增加計數器
	// }

	// ticker := time.Tick(1000 * time.Millisecond)

	// for range ticker {
	//     func() {
	//         wg.Add(1)
	//         defer wg.Done()

	//         doObj.Do_choose_api("DO_WHOLE", do, 0, "")
	// 		res := doObj.Do_choose_api("DO_CHECK", do, 0, "1")

	//         if res == lastCheckResult {
	//             return
	//         }
	//         lastCheckResult = res
	//         go task1(doObj, do)
	//     }()
	// }
	// go func() {
	// 	defer wg.Done()
	// 	task2(doObj, di)
	// }()

	// wg.Wait() // 等待所有Goroutines完成

	input := "device{(e1200, 192.168.127.254)}, when{[(do,0)=0]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]]}"
	extracted := extractContent(input)

	input2 := "device{(e1200, 192.168.127.254)}, when{[(do,2)=1]}, then{[(do,4)=1]&&[(do,5)=1]&&[(do,6)=1]}"
	extracted2 := extractContent(input2)

	taskQueue.AddTask(func(taskInfo Task) {
		task_func(extracted)
	})
	taskQueue.AddTask(func(taskInfo Task) {
		task_func(extracted2)
	})
	go func() {
		for {
			taskQueue.ExecuteTasks()
			time.Sleep(time.Second)
		}
	}()
	select {}
	fmt.Println("All Goroutines have finished.")

}
func extractContent(input string) Task {
	var task Task

	// device
	deviceRe := regexp.MustCompile(`device{\(([^,]+),\s*([^}]+)\)}`)
	deviceMatch := deviceRe.FindStringSubmatch(input)
	if len(deviceMatch) == 3 {
		task.Device = DeviceInfo{
			Type: deviceMatch[1],
			IP:   deviceMatch[2],
		}
	}

	// when
	whenRe := regexp.MustCompile(`when{\[([^\]]+)\]}`)
	whenMatch := whenRe.FindStringSubmatch(input)

	if len(whenMatch) == 2 {
		condition := strings.Trim(whenMatch[1], "[]")
		parts := strings.Split(strings.Trim(condition, "[]"), "=")

		if len(parts) == 2 {
			// actionParts := strings.Split(parts[0], ",")
			re := regexp.MustCompile(`\(([^)]+)\)`)
			chParts := re.FindStringSubmatch(parts[0])

			if len(chParts) == 2 {
				chTypeNum := strings.Split(chParts[1], ",")
				task.When = append(task.When, Condition{
					ChType: strings.TrimSpace(chTypeNum[0]),
					ChNum:  strings.TrimSpace(chTypeNum[1]),
					Value:  strings.TrimSpace(parts[1]),
				})
			}
		}
	}
	//then
	thenRe := regexp.MustCompile(`then{([^}]+)}`)
	thenMatch := thenRe.FindStringSubmatch(input)

	if len(thenMatch) == 2 {
		actions := strings.Split(thenMatch[1], "&&")

		for _, action := range actions {
			parts := strings.Split(strings.Trim(action, "[]"), "=")
			if len(parts) == 2 {
				// actionParts := strings.Split(parts[0], ",")
				re := regexp.MustCompile(`\(([^)]+)\)`)
				chParts := re.FindStringSubmatch(parts[0])

				if len(chParts) == 2 {
					chTypeNum := strings.Split(chParts[1], ",")
					task.ThenActions = append(task.ThenActions, Action{
						ChType: strings.TrimSpace(chTypeNum[0]),
						ChNum:  strings.TrimSpace(chTypeNum[1]),
						Value:  strings.TrimSpace(parts[1]),
					})
				}
			}
		}
	}

	return task
}
func task_func(taskInfo Task) {
	// var lastCheckResult int
	doObj := do.DoObj{}
	do := do.NewMachine(taskInfo.Device.Type, "1213", taskInfo.Device.IP, "do", 8)
	// doObj.Do_choose_api("DO_WHOLE", do, 0, "")
	res := 1
	// for _, condi := range taskInfo.When { // 遞迴比較條件，皆為1才會過
	//     chnum, _ := strconv.Atoi(condi.ChNum)
	//     res &= doObj.Do_choose_api("DO_CHECK", do, chnum, condi.Value)
	// }

	// if res == taskInfo.LastCheck {
	// 	// check last
	// 	fmt.Println("Status is the same as last check, no action needed.")
	// 	return
	// }
	// taskInfo.LastCheck = res // update status
	if res == 1 {
		for _, action := range taskInfo.ThenActions {
			chnum, _ := strconv.Atoi(action.ChNum)
			doObj.Do_choose_api("TEST_DO_PUT_VALUE", do, chnum, action.Value)
		}
	}
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue: make([]func(task Task), 0),
	}
}

func (q *TaskQueue) AddTask(taskFunc func(task Task)) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, taskFunc)
}

func (q *TaskQueue) ExecuteTasks() {
	for _, taskFunc := range q.queue {
		taskFunc(Task{})
	}
}

func task1(doObj do.DoObj, do *do.Machine) {

	res := doObj.Do_choose_api("DO_CHECK", do, 0, "1")
	if res == 1 {
		sub_task1(doObj, do)
	} else {
		sub_task2(doObj, do)
	}

	fmt.Println("End task1")
}

func sub_task1(doObj do.DoObj, do *do.Machine) {

	doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "1")
	fmt.Println("End sub task1")
	// do.Do_clear_ch(do, 0)
}

func sub_task2(doObj do.DoObj, do *do.Machine) {

	doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "0")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "0")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "0")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "0")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "0")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 6, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 7, "1")
	fmt.Println("End sub task2")
	// do.Do_clear_ch(do, 0)
}
