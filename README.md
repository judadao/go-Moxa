# Moxa resful control module

# Feature

- Control Moxa devices through a RESTful API.
  
  - Divide the RESTful API of each Moxa device into controllable functions, where operating a function only requires inputting the corresponding keyword.
    
  - Using object declaration for devices makes usage more intuitive.
    
  - Freedom to assemble function requirements.
  - (2024/05/21 update) A task can be automatically generated from the input sentence task
  - (#TODO: 串接AI)
    
- Cross-platform control supports the following devices:
  
  - Moxa E1200 Series(DO)
    

## Quick Start

### Control machines with sentences
#### Excute 
- Install
  1. Clone this repository: git clone [repo_url]
  2. Run:
    ```bash
    go run main.go
    ```
#### Prompt
- Create Your Own Prompt
  - Objects will be created using "( )"
    ex1.
    ```go
    //In When & Then 
    (device1,do,0) // name:device1, channel type: do, channel nmber: 0  
    ```
    ex2.
    ```go
    //In Device
    (e1200, e1213, 192.168.127.254) // type:e1200, sub type: e1213, IP: 192.168.127.254
    ```
  - Conditions will be set using  "[ ]"
    ex1.
    ```go
    //In When & Then 
    [(device1, do,0)=1] 
    // When: device1's do channel 0 equals 1
    // Then: device1's do channel 0 is set to 1
    ```
    ex2.
    ```go
    //In Device
    [(e1200, e1213, 192.168.127.254)=device1] // set device1 info
    ```
  - Condition Groups will use "{ }" in three main fields
  
    ```go  
    device{}, when{}, then{}
    ```
  - Example of a Full Prompt
    ```go
    device{[(e1200, e1213, 192.168.127.254)=device1]&&[(e1200, e1213, 192.168.127.255)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,1)=0]}, then{[(device1,do,2)=1]&&[(device1, do,3)=1]}
    // device1 = (e1200, e1213, 192.168.127.254)
    // device2 = (e1200, e1213, 192.168.127.255)
    // When device1's do channel 0 is equal to 1 and device2's do channel 1 is equal to 0 
    // Then device1's do channel 2 is set to 1 and device1's do channel 3 is set to 1
    ```
- Add an Task
  ```go
    input1 := "device{(e1200, e1213, 192.168.127.254)}, when{[(do,0)=1]}, then{[(do,1)=0]&&[(do,2)=0]&&[(do,3)=0]}"
    extracted1 := task_qu.ExtractContent(input1) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted1)
    })
  ```
- Add many Tasks
  ```go
    input1 := "device{[(e1200, e1213, 192.168.127.254)=device1]&&[(e1200, e1213, 192.168.127.254)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,1)=0]}, then{[(device1,do,2)=1]&&[(device1, do,3)=1]}"
    extracted1 := task_qu.ExtractContent(input1) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted1) // 將info輸入task
    })

    input2 := "device{[(e1200, e1213, 192.168.127.254)=device1]&&[(e1200, e1213, 192.168.127.255)=device2]}, when{[(device1, do,0)=1]&&[(device2, do,1)=1]}, then{[(device1,do,2)=0]&&[(device1, do,3)=0]}"
    extracted2 := task_qu.ExtractContent(input2) //將input轉換成info

    taskQueue.AddTask(func(taskInfo task_qu.Task) { //將task加入queue
      allTask.Task_func(&extracted2)
    })
  ```



### Assemble everything you need

- Declare your Device as an object
  
  ```go
    doObj := do.NewMachine("e1200", "1213","192.168.127.254", "do", 8)
    //Parameters are "main model", "sub-model","IP", "IO channel type", "Channel numbers"
  ```
  

  
- Use Do_choose_api to select the RESTful function you want to execute. Below are examples using get / put Value.
  
- Get:
  
  ```go
    do.Do_choose_api("DO_GET_VALUE", doObj, 1, "")
    //Parameters: "function keyword", "machine obj", "channel number"
  ```
  
- Put:
  
  ```go
    do.Do_choose_api("DO_PUT_VALUE", doObj, 1, "0")
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
    
                do.Do_choose_api("DO_WHOLE", doObj, 0, "") //update resful value
                res := do.Do_choose_api("DO_CHECK", doObj, 0, "1")//check machine channel 0 status
    
                if res == lastCheckResult { ////ensure only executes once
                    return
                }
                lastCheckResult = res
                go task1(doObj) //execute task1
            }()
        }
    ```
  
- task1
  
  ```go
    func task1( doObj *do.Machine) {
    
        res :=do.Do_choose_api("DO_CHECK", doObj, 0, "1") //check if ch is 1
        if res == 1 {
            sub_task1(doObj) // true execute sub task1
        }else{
            sub_task2(doObj) // true execute sub task2
        }
    
        
        fmt.Println("End task1")
    }
    ```
  
- sub task
  
  ```go
    func sub_task1( doObj *do.Machine) { //set do ch 1~5 to on
    
        do.Do_choose_api("DO_PUT_VALUE", doObj, 1, "1")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 2, "1")
        do.Do_choose_api("DO_PUT_VALUE", do, 3, "1")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 4, "1")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 5, "1")
        fmt.Println("End sub task1")
    }
    
    func sub_task2(doObj *do.Machine) { //set do ch 1~5 to off, 6~7 to on
    
        do.Do_choose_api("DO_PUT_VALUE", doObj, 1, "0")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 2, "0")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 3, "0")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 4, "0")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 5, "0")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 6, "1")
        do.Do_choose_api("DO_PUT_VALUE", doObj, 7, "1")
        fmt.Println("End sub task2")
    }
    ```
