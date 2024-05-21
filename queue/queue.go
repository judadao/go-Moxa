package task_qu

import (
	"regexp"
	"strings"
	"sync"
)

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
	SubType string
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

func ExtractContent(input string) Task {
	var task Task

	// device
	deviceRe := regexp.MustCompile(`device{\(([^,]+),\s*([^}]+),\s*([^}]+)\)}`)
	deviceMatch := deviceRe.FindStringSubmatch(input)
	if len(deviceMatch) == 4 {
		task.Device = DeviceInfo{
			Type: deviceMatch[1],
			SubType: deviceMatch[2],
			IP:   deviceMatch[3],
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
		go taskFunc(Task{})
	}
}