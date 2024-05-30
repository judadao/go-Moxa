package main

import (
	"fmt"
	task_qu "go-Moxa/queue"
	allTask "go-Moxa/task"
	"time"
)

// var globalDoObjects []*do.Machine



func main() {
	


	taskQueue := task_qu.NewTaskQueue()


	input2 := "device{[(1200, e1213, 192.168.127.254)=device1]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,0)=1]}, then{[(device1,do,2)=0]&&[(device1, do,3)=0]}"
	extracted2 := task_qu.ExtractContent(input2)
	// fmt.Println(extracted2.When)

	input := "device{[(1200, e1213, 192.168.127.254)=device1]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,0)=0]}, then{[(device1,do,2)=1]&&[(device1, do,3)=1]}"
	extracted := task_qu.ExtractContent(input)

	


	taskQueue.AddTask(func(taskInfo task_qu.Task) {
		allTask.Task_func(&extracted2)
	})


	taskQueue.AddTask(func(taskInfo task_qu.Task) {
		allTask.Task_func(&extracted)
	})
	go func() {
		for {
			taskQueue.ExecuteTasks()
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
	}()
	select {}
	fmt.Println("All Goroutines have finished.")

	// do4510 := do.NewMachine_4510("4510", "2600", "192.168.127.254", "45MR-2600-0", "do", 16)
	// do1200 := do.NewMachine("1200", "e1213", "192.168.127.253", "do",8)
	// do.Do_choose_api("DO_GET_VALUE", do1200, 0, "0")
	// do.Do_choose_api("DO_GET_VALUE", do4510, "DO-00", "0")
	// do.Do_choose_api("DO_GET_VALUE", do4510, "DO-01", "0")

}




