package main

import (
	"fmt"
	"go-Moxa/do"
	"sync"
	"time"
)
var wg sync.WaitGroup
func main() {
	// router = chi.NewRouter()
	
	// concurrency := 100 //限制100條thread
	var lastCheckResult int
	
	doObj := do.DoObj{}
	do := do.NewMachine("e1200", "1213","192.168.127.254", "do", 8)
	// di2 := do.NewMachine("e1200", "1213","192.1.1.2", "di", 8)
	// di3 := do.NewMachine("e1200", "1213","192.1.1.3", "di", 8)
	// di4 := do.NewMachine("e1200", "1213","192.1.1.4", "di", 8)
	// router.Get("/", doObj.Do_choose_api)
	
	// for i := 0; i < 4; i++ {
	// 	wg.Add(1) // 增加計數器
	// }

	
	ticker := time.Tick(1000 * time.Millisecond) 

    for range ticker {
        func() {
            wg.Add(1)
            defer wg.Done()

            doObj.Do_choose_api("DO_WHOLE", do, 0, "")
			res := doObj.Do_choose_api("DO_CHECK", do, 0, "1")

            if res == lastCheckResult {
                return
            }
            lastCheckResult = res
            go task1(doObj, do)
        }()
    }
	// go func() {
	// 	defer wg.Done()
	// 	task2(doObj, di)
	// }()

	
	//To 看現有schedual，來測試schel task
	//併發執行完的
	
	wg.Wait() // 等待所有Goroutines完成
 	fmt.Println("All Goroutines have finished.")

}

func task1(doObj do.DoObj, do *do.Machine) {

	res :=doObj.Do_choose_api("DO_CHECK", do, 0, "1")
	if res == 1 {
		// wg.Add(1)
		// go func() {
			// defer wg.Done()
			sub_task1(doObj, do)
			// doObj.Do_choose_api("DO_PUT_VALUE", do, 0, "1")
		// }()
	}else{
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
			sub_task2(doObj, do)
			// doObj.Do_choose_api("DO_PUT_VALUE", do, 0, "0")

		// }()
	}

	
	//##TODO: 完成put api
	
	// doObj.Do_choose_api("DO_PAULSESTATUS", di, 0)

	// doObj.Do_choose_api("DO_READ_CH", di, 0)
	// doObj.Do_choose_api("DO_READ_CH", di, 0)
	// fmt.Println( do.Do_pop_ch(di, 0))
	

	
	fmt.Println("End task1")
}

func sub_task1(doObj do.DoObj, do *do.Machine) {

	doObj.Do_choose_api("DO_PUT_VALUE", do, 1, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 2, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 3, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 4, "1")
	doObj.Do_choose_api("DO_PUT_VALUE", do, 5, "1")
	fmt.Println("End sub task1")
	// do.Do_clear_ch(do, 0)
}

func sub_task2(doObj do.DoObj, do *do.Machine) {

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