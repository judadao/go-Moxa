package main

import (
	"fmt"
	"go-Moxa/do"
	"sync"
)

func main() {
	// router = chi.NewRouter()
	var wg sync.WaitGroup
	// concurrency := 100 //限制100條thread
	
	doObj := do.DoObj{}
	di := do.NewMachine("e1200", "1213","192.168.127.254", "do", 8)
	// di2 := do.NewMachine("e1200", "1213","192.1.1.2", "di", 8)
	// di3 := do.NewMachine("e1200", "1213","192.1.1.3", "di", 8)
	// di4 := do.NewMachine("e1200", "1213","192.1.1.4", "di", 8)
	// router.Get("/", doObj.Do_choose_api)
	
	// for i := 0; i < 4; i++ {
	// 	wg.Add(1) // 增加計數器
	// }
	wg.Add(1)
	go func() {
		defer wg.Done()
		task1(doObj, di)
	}()
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
	// defer close(di.Channel[0])
	// doObj.Do_choose_api("DO_WHOLE", di, 0)
	// fmt.Println( do.Do_pop_ch(di, 0))
	
	// doObj.Do_choose_api("DO_STATUS", di, 0)
	doObj.Do_choose_api("DO_GET_VALUE", do, 1)
	doObj.Do_choose_api("DO_GET_VALUE", do, 2)
	doObj.Do_choose_api("DO_GET_VALUE", do, 3)
	doObj.Do_choose_api("DO_GET_VALUE", do, 4)
	doObj.Do_choose_api("DO_GET_VALUE", do, 5)
	//##TODO: 完成put api
	
	// doObj.Do_choose_api("DO_PAULSESTATUS", di, 0)

	// doObj.Do_choose_api("DO_READ_CH", di, 0)
	// doObj.Do_choose_api("DO_READ_CH", di, 0)
	// fmt.Println( do.Do_pop_ch(di, 0))
	

	
	fmt.Println("End task1")
}

func task2(doObj do.DoObj, di *do.Machine) {
	doObj.Do_choose_api("DO_PAULSESTATUS", di, 0)
	// do.Do_pop_ch(di, 0)
	// doObj.Do_choose_api("DO_PAULSESTATUS", di, 0)
	// doObj.Do_choose_api("DO_PAULSECOUNT", di, 0)
	fmt.Println("End task2")
	do.Do_clear_ch(di, 0)
}