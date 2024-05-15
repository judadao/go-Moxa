# Moxa resful control module

# Feature

- 控制Moxa device 透過restful api
  
  - 將Moxa 各device 的restful api切分成可操控的function，且操作function只需要填入對應的關鍵子即可使用
    
  - Device 使用物件宣告使用上更直觀
    
  - 可以自由拼裝Function需求
    
- 可跨平台控制可支援的裝置如下
  
  - moxa E1200 系列
    
    - DO
      
    - (#TODO) DI
      
  - (#TODO) moxa 4510
    
  - (#TODO) moxa 後續新機種
    

## Quick Start

- 宣告你的Device為object
  
  ```go
    do := do.NewMachine("e1200", "1213","192.168.127.254", "do", 8)
    //參數分別為"主型號", "副型號","IP", "IO channel type", "Channel numbers"
  ```
  
- 宣告Do interface
  
  ```go
    doObj := do.DoObj{}
  ```
  
- 使用Do_choose_api選擇你要執行的restful功能。以下用get /put Value來舉例
  
- Get:
  
  ```go
    doObj.Do_choose_api("DO_GET_VALUE", do, 1, "")
    //參數"function關鍵字", "machine obj", "channel number"
  ```
  
- Put:
  
  ```go
    doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "0")
    //參數"function關鍵字", "machine obj", "channel number", "msg"
    ```
  
- 目前可支援E1200
  
  - DO
    
    - DO_WHOLE : update 目前本機儲存的restful value
      
    - DO_CHECK : 檢查目前的channel值
      
    - DO_PUT_VALUE : 將value put進對應的channel
      
    - DO_GET_VALUE : 取得對應Channel的值
      

## You can do

- 在main.go中有範例
  
  - 展示用update, check 跟put 的組合，完成使用go moxa控制1200 do
    
- 首先我希望每一秒去檢查目前的restful 狀態，符合條件後觸發task1
  
     ```go
    ticker := time.Tick(1000 * time.Millisecond)  //1 sec 
    
        for range ticker {
            func() {
                wg.Add(1)
                defer wg.Done()
    
                doObj.Do_choose_api("DO_WHOLE", do, 0, "") //update resful value
                res := doObj.Do_choose_api("DO_CHECK", do, 0, "1")//check machine channel 0 status是否是1
    
                if res == lastCheckResult { //確保只執行一次不會一直執行
                    return
                }
                lastCheckResult = res
                go task1(doObj, do) //執行task1
            }()
        }
    ```
  
- task1
  
  ```go
    func task1(doObj do.DoObj, do *do.Machine) {
    
        res :=doObj.Do_choose_api("DO_CHECK", do, 0, "1") //判斷ch是否為1
        if res == 1 {
            sub_task1(doObj, do) // true 執行sub task1
        }else{
            sub_task2(doObj, do) // false 執行sub task2
        }
    
        
        fmt.Println("End task1")
    }
    ```
  
- sub task
  
  ```go
    func sub_task1(doObj do.DoObj, do *do.Machine) { //do ch 1~5 設為on
    
        doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "1")
        doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "1")
        fmt.Println("End sub task1")
        // do.Do_clear_ch(do, 0)
    }
    
    func sub_task2(doObj do.DoObj, do *do.Machine) { // do ch 1~5設為off, 6~7 設為on
    
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
    ```
