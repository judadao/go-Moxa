package main

import (
	task_qu "go-Moxa/queue"
	allTask "go-Moxa/task"
	"time"
)

// var globalDoObjects []*do.Machine



func main() {
	


	taskQueue := task_qu.NewTaskQueue()

	input := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=0]}, then{[(do,1)=1]&&[(do,2)=1]&&[(do,3)=1]&&[(do,4)=0]&&[(do,5)=0]&&[(do,6)=0]}"
	extracted := task_qu.ExtractContent(input)

	input2 := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=1]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]&&[(do,4)=1]&&[(do,5)=1]&&[(do,6)=1]}"
	extracted2 := task_qu.ExtractContent(input2)

	taskQueue.AddTask(func(taskInfo task_qu.Task) {
		allTask.Task_func(&extracted)
	})
	taskQueue.AddTask(func(taskInfo task_qu.Task) {
		allTask.Task_func(&extracted2)
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