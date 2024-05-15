# Moxa resful control module

# Feature

- Control Moxa devices through a RESTful API.
  
  - Divide the RESTful API of each Moxa device into controllable functions, where operating a function only requires inputting the corresponding keyword.
    
  - Using object declaration for devices makes usage more intuitive.
    
  - Freedom to assemble function requirements.
    
- Cross-platform control supports the following devices:
  
  - Moxa E1200 Series
    
    - DO
      
    - (#TODO) DI
      
  - (#TODO) moxa 4510
    
  - (#TODO) Future Moxa models
    

## Quick Start

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
      

## You can do

- There is an example in main.go
  
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
