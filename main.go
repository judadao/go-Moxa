package main

import (
	task_qu "go-Moxa/queue"
	allTask "go-Moxa/task"
	"time"
)

// var globalDoObjects []*do.Machine



func main() {
	


	taskQueue := task_qu.NewTaskQueue()


	input2 := "device{[(e1200, e1213, 192.168.127.254)=device1]&&[(e1200, e1213, 192.168.127.254)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,1)=1]}, then{[(device1,do,2)=0]&&[(device1, do,3)=0]}"
	extracted2 := task_qu.ExtractContent(input2)

	input := "device{[(e1200, e1213, 192.168.127.254)=device1]&&[(e1200, e1213, 192.168.127.254)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,1)=0]}, then{[(device1,do,2)=1]&&[(device1, do,3)=1]}"
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
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()
	select {}
	// fmt.Println("All Goroutines have finished.")

}