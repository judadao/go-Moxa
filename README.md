# Moxa resful control module

# Feature

- Control Moxa devices through a RESTful API.
  
  - Divide the RESTful API of each Moxa device into controllable functions, where operating a function only requires inputting the corresponding keyword.
    
  - Using object declaration for devices makes usage more intuitive.
    
  - Freedom to assemble function requirements.
  - (2024/05/21 update) A task can be automatically generated from the input sentence.task
  - (#TODO: 串接AI)
    
- Cross-platform control supports the following devices:
  
  - Moxa E1200 Series(DO)
    

## Quick Start

### Control machines with sentences

(#TODO: 跨Device 之間的prompt/struct/parse)
- 製作你的prompt
  ```go
    "device{(type , sub_type, IP)}, when{ condition}, then{action}"
  ```
  ex. 
  ```go
    "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=1]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]}"
    //Device: e1200, sutype: e1213, IP:192.168.127.254
    //when Do channel 0 是 1 的時候
    //Then Do Channel 1, 2, 3設為0
  ```
- 加入你的Task
  ```go
    input1 := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=1]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]}"
    extracted1 := task_qu.ExtractContent(input1) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted1)
    })
  ```
- 多個Task
  ```go
    input1 := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=1]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]}"
    extracted1 := task_qu.ExtractContent(input1) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted1) // 將info輸入task
    })

    input2 := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=0]}, then{[(do,1)=1]&&[(do,2)=1]&&[(do,3)=1]}"
    extracted2 := task_qu.ExtractContent(input2) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted2)
    })
  ```



### Assemble everything you need

- Declare your Device as an object
  
  ```go
    do := do.NewMachine("e1200", "1213","192.168.127.254", "do", 8)
    //Parameters are "main model", "sub-model","IP", "IO channel type", "Channel numbers"
  ```
  
- Declare Do interface
  
  ```go
    doObj := do.DoObj{}
  ```
  
- Use Do_choose_api to select the RESTful function you want to execute. Below are examples using get / put Value.
  
- Get:
  
  ```go
    doObj.Do_choose_api("DO_GET_VALUE", do, 1, "")
    //Parameters: "function keyword", "machine obj", "channel number"
  ```
  
- Put:
  
  ```go
    doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "0")
    //Parameters: "function keyword", "machine obj", "channel number", "msg"
    ```
  
- Currently supports E1200
  
  - DO
    
    - DO_WHOLE: Update the current locally stored RESTful value
      
    - DO_CHECK: Check the current channel value
      
    - DO_PUT_VALUE: Put the value into the corresponding channel
      
    - DO_GET_VALUE: Get the value of the corresponding channel
      

## Demonstrate the usage of classes and functions

- You can freely assemble functions and objects as needed.
- Demonstrates the combination of update, check, and put to control the E1200 DO using Go Moxa.
    
- First, I want to check the current RESTful status every second and trigger task1 when conditions are met.
  
     ```go
    ticker := time.Tick(1000 * time.Millisecond)  //1 sec 
    
        for range ticker {
            func() {
                wg.Add(1)
                defer wg.Done()
    
                doObj.Do_choose_api("DO_WHOLE", do, 0, "") //update resful value
                res := doObj.Do_choose_api("DO_CHECK", do, 0, "1")//check machine channel 0 status
    
                if res == lastCheckResult { ////ensure only executes once
                    return
                }
                lastCheckResult = res
                go task1(doObj, do) //execute task1
            }()
        }
    ```
  
- task1
  
  ```go
    func task1(doObj do.DoObj, do *do.Machine) {
    
        res :=doObj.Do_choose_api("DO_CHECK", do, 0, "1") //check if ch is 1
        if res == 1 {
            sub_task1(doObj, do) // true execute sub task1
        }else{
            sub_task2(doObj, do) // true execute sub task2
        }
    
        
        fmt.Println("End task1")
    }
    ```
  
- sub task
  
  ```go
    func sub_task1(doObj do.DoObj, do *do.Machine) { //set do ch 1~5 to on
    
        doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "1")
        fmt.Println("End sub task1")
    }
    
    func sub_task2(doObj do.DoObj, do *do.Machine) { //set do ch 1~5 to off, 6~7 to on
    
        doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "0")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "0")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "0")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "0")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "0")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 6, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 7, "1")
        fmt.Println("End sub task2")
    }
    ```
