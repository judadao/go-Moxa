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
	Device       map[string]DeviceInfo
	When        []Condition
	ThenActions []Action
	LastCheck   int //check last status
}

type DeviceInfo struct {
	Type string
	SubType string
	IP   string
	Name    string
}

type Condition struct {
	ObjName string
	ChType string
	ChNum  string
	Value  string
}

type Action struct {
	ObjName string
	ChType string
	ChNum  string
	Value  string
}

func ExtractContent(input string) Task {
	// var task Task
	task := Task{Device: make(map[string]DeviceInfo)}
	task.LastCheck = -1

	// device
	deviceRe := regexp.MustCompile(`\[\(([^,]+),\s*([^,]+),\s*([^\)]+)\)\s*=\s*([^\]]+)\]`)
	deviceMatches := deviceRe.FindAllStringSubmatch(input, -1)
	for _, match := range deviceMatches {
		if len(match) == 5 {
			device := DeviceInfo{
				Type:    match[1],
				SubType: match[2],
				IP:      match[3],
				Name:    match[4],
			}
			task.Device[device.Name] = device
		}
	}

	// when
	whenRe := regexp.MustCompile(`when{([^}]+)}`)
	whenMatch := whenRe.FindStringSubmatch(input)
	if len(whenMatch) == 2 {
		conditions := strings.Split(whenMatch[1], "&&")

		for _, condition := range conditions {
			parts := strings.Split(strings.Trim(condition, "[]"), "=")
			
			if len(parts) == 2 {
				// actionParts := strings.Split(parts[0], ",")
				re := regexp.MustCompile(`\(([^)]+)\)`)
				chParts := re.FindStringSubmatch(parts[0])

				if len(chParts) == 2 {
					chTypeNum := strings.Split(chParts[1], ",")
					task.When = append(task.When, Condition{
						ObjName: strings.TrimSpace(chTypeNum[0]),
						ChType: strings.TrimSpace(chTypeNum[1]),
						ChNum:  strings.TrimSpace(chTypeNum[2]),
						Value:  strings.TrimSpace(parts[1]),
					})
				}
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
						ObjName: strings.TrimSpace(chTypeNum[0]),
						ChType: strings.TrimSpace(chTypeNum[1]),
						ChNum:  strings.TrimSpace(chTypeNum[2]),
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