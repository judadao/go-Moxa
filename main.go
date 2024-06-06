package main

import (
	"bufio"
	"context"
	"fmt"
	task_qu "go-Moxa/queue"
	allTask "go-Moxa/task"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// var globalDoObjects []*do.Machine



func main() {
	
	key := "key"

	client := openai.NewClient(key)

	taskQueue := task_qu.NewTaskQueue()

	//4510: type, subtype, ip, slot name
	//1200: type, subtype, ip, slot name(slot name no need to set)
	// input2 := "device{[(4510, 2600, 192.168.127.254)=45MR-2600-0]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(45MR-2600-0, do,DO-00)=1]&&[(device2, do,0)=1}, then{[(45MR-2600-0,do,DO-01)=1]&&[(45MR-2600-0, do,DO-02)=1]}"
	// extracted2 := task_qu.ExtractContent(input2)

	// input1 := "device{[(4510, 2600, 192.168.127.254)=45MR-2600-0]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(45MR-2600-0, do,DO-00)=0]&&[(device2, do,0)=1}, then{[(45MR-2600-0,do,DO-01)=0]&&[(45MR-2600-0, do,DO-02)=0]}"
	// extracted1 := task_qu.ExtractContent(input1)
	// fmt.Println(extracted2.When)

	// input := "device{[(1200, e1213, 192.168.127.254)=device1]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,0)=0]}, then{[(device1,do,2)=1]&&[(device1, do,3)=1]}"
	// extracted := task_qu.ExtractContent(input)

	


	// taskQueue.AddTask(func(taskInfo task_qu.Task) {
	// 	allTask.Task_func(&extracted2)
	// })


	// taskQueue.AddTask(func(taskInfo task_qu.Task) {
	// 	allTask.Task_func(&extracted1)
	// })

	go startConversation(client, taskQueue)
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



// processInput 將輸入的字符串切分並拼接成帶有 "&&" 的格式
func processInput(input string) string {
	entries := strings.Split(input, "], ")
	for i := 0; i < len(entries); i++ {
		entries[i] = strings.TrimSpace(entries[i])
		if !strings.HasSuffix(entries[i], "]") {
			entries[i] += "]"
		}
	}
	return strings.Join(entries, "&&")
}

func assembleInfo(deviceInfo, triggerConditions, actions []string) string {
	deviceStr := "device{" + strings.Join(deviceInfo, "&&") + "}"
	triggerStr := "when{" + strings.Join(triggerConditions, "&&") + "}"
	actionStr := "then{" + strings.Join(actions, "&&") + "}"

	return fmt.Sprintf("%s, %s, %s", deviceStr, triggerStr, actionStr)
}

//TODO: add case:放棄set device/顯示目前有的任務 
func startConversation(client *openai.Client, task_queue *task_qu.TaskQueue) {
		// 創建一個用於讀取用戶輸入的Scanner
		scanner := bufio.NewScanner(os.Stdin)

		// 創建一個用於記錄用戶輸入的slice
		var conversation []openai.ChatCompletionMessage
		var deviceInfo []string
		var triggerConditions []string
		var actions []string
		var extrDict = make(map[string] *task_qu.Task)

		
		for {
			// 讀取用戶輸入
			// 初始提示信息
			fmt.Println("請輸入 'moxarest' 來開始設置，或是進行一般對話")
			fmt.Print("You: ")
			scanner.Scan()
			userInput := scanner.Text()
			shouldBreak := false
			// userInput :="moxarest"

			if strings.ToLower(userInput) == "moxarest" {
				
				conversation = append(conversation, openai.ChatCompletionMessage{Role: "system", Content: "請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)"})
				fmt.Println(conversation[len(conversation)-1].Content)

				for !shouldBreak {
				
					fmt.Print("You: ")
					scanner.Scan()
					userInput = scanner.Text()
					conversation = append(conversation, openai.ChatCompletionMessage{Role: "user", Content: userInput})

					// 根據用戶輸入添加不同的對話內容
					switch userInput {
					case "機器":
					case "1":
						fmt.Println("MOXAREST: 請提供裝置信息，格式為[(主類型，子類型，IP地址)=裝置名稱]，多個裝置信息之間用逗號分隔：")
						fmt.Print("You: ")
						// scanner.Scan()
						// deviceInput := scanner.Text()
						deviceInput :="[(4510, 2600, 192.168.127.254)=45MR-2600-0], [(1200, e1211, 192.168.127.253)=device2]"
						
						joinedDeviceInfo := processInput(deviceInput)

						deviceInfo = append(deviceInfo, joinedDeviceInfo)
						conversation = append(conversation, openai.ChatCompletionMessage{Role: "user", Content: joinedDeviceInfo})

						
						deviceString := strings.Join(deviceInfo, "&&")
						fmt.Printf("device{%s}\n", deviceString)
						fmt.Println("請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)")
					case "觸發條件":
					case "2":
						fmt.Println("MOXAREST: 請提供觸發條件信息，格式為[(裝置名稱，通道類型，通道名稱)=值]，多個觸發條件之間用逗號分隔：")
						fmt.Print("You: ")
						scanner.Scan()
						triggerInput := scanner.Text()
						// triggerInput := "[(45MR-2600-0, do, DO-00)=1], [(device2, do, 0)=1]"
						joinedTriggerConditions := processInput(triggerInput)
						triggerConditions = append(triggerConditions, joinedTriggerConditions)
						conversation = append(conversation, openai.ChatCompletionMessage{Role: "user", Content: joinedTriggerConditions})
						fmt.Println("請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)")

					case "動作":
					case "3":
						fmt.Println("MOXAREST: 請提供動作信息，格式為[(裝置名稱，通道類型，通道名稱)=值]，多個動作之間用逗號分隔：")
						fmt.Print("You: ")
						scanner.Scan()
						actionInput := scanner.Text()
						// actionInput := "[(45MR-2600-0, do, DO-01)=0], [(45MR-2600-0, do, DO-02)=0]"
						joinedActions := processInput(actionInput)

						actions = append(actions, joinedActions)
						conversation = append(conversation, openai.ChatCompletionMessage{Role: "user", Content: joinedActions})
						fmt.Println("請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)")

					case "顯示":
					case "4":
						// 顯示目前的 device、when 和 action 資訊
						deviceString := strings.Join(deviceInfo, "&&")
						triggerString := strings.Join(triggerConditions, "&&")
						actionString := strings.Join(actions, "&&")
						fmt.Printf("device{%s}\n", deviceString)
						fmt.Printf("when{%s}\n", triggerString)
						fmt.Printf("action{%s}\n", actionString)
						fmt.Println("請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)")

					case "End":
					case "5":
						// 結束設置過程
						fmt.Println("設定結束。")
						fmt.Println("MOXAREST: 此Task的名字是?")
						fmt.Print("You: ")
						scanner.Scan()
						endInput := scanner.Text()
						result := assembleInfo(deviceInfo, triggerConditions, actions)
						fmt.Println(result)

						// addResultToGlobalDict(endInput, result)
						extrDict[endInput] = task_qu.ExtractContent(result)
						// extractedResult := task_qu.ExtractContent(result)
						
						// extractedResult2 := task_qu.ExtractContent("device{[(4510, 2600, 192.168.127.254)=45MR-2600-0]&&[(1200, e1211, 192.168.127.253)=device2]}, when{[(45MR-2600-0, do,DO-00)=0]&&[(device2, do,0)=1]}, then{[(45MR-2600-0,do,DO-01)=1]&&[(45MR-2600-0, do,DO-02)=1]}")

						task_queue.AddTask(func(taskInfo task_qu.Task) {
							allTask.Task_func(extrDict[endInput])
						})

						// task_queue.AddTask(func(taskInfo task_qu.Task) {
						// 	allTask.Task_func(&extractedResult2)
						// })

						deviceInfo = nil
    					triggerConditions = nil
    					actions = nil
						// fmt.Println("請輸入 'moxarest' 來開始設置，或是進行一般對話")
						shouldBreak = true
						

					default:
						// 發送對話內容到OpenAI API
						fmt.Println("關鍵詞有誤，請輸入以下其中一個關鍵詞：機器(1)、觸發條件(2)、動作(3)、顯示(4)、結束(5)")
						fmt.Println("或是ENDˋ結束設定進行一般對話")
					}
				}
			} else {
				// 發送對話內容到OpenAI API
				conversation = append(conversation, openai.ChatCompletionMessage{Role: "user", Content: userInput})
				resp, err := client.CreateChatCompletion(
					context.Background(),
					openai.ChatCompletionRequest{
						Model:    openai.GPT3Dot5Turbo,
						Messages: conversation,
					},
				)
				if err != nil {
					fmt.Printf("ChatCompletion error: %v\n", err)
					break
				}

				// 打印OpenAI API返回的回復
				fmt.Println("AI:", resp.Choices[0].Message.Content)

				// 將AI回復添加到對話內容中
				conversation = append(conversation, openai.ChatCompletionMessage{Role: "assistant", Content: resp.Choices[0].Message.Content})

				// fmt.Println("請輸入 'moxarest' 來開始設置，或是進行一般對話")
			}
		}
	
}


