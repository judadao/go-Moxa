package do

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type do_allApiFunc func(string)

var doApiMap = map[string]do_allApiFunc{
	"DO_WHOLE": do_get_whole_value,
	"DO_STATUS": do_get_status,
	"DO_PAULSESTATUS": do_get_paulse_status,
	"DO_PAULSECOUNT": do_get_paulse_count,

}

func muti_thread_api(callback func()){
	callback();
}


func Do_choose_api(w http.ResponseWriter, r *http.Request) {
	var apiKey = ""
	fmt.Println("yes this is do")
	w.Write([]byte("Hello, world!"))
	apiKey = "DO_WHOLE"
	if fn, ok := doApiMap[apiKey];ok{
		fn("i'm "+apiKey)
	}else{
		fmt.Println("not exit")
	}
	
	
}

func do_get_whole_value(msg string){
	fmt.Println("do_get_whole_value:", msg)
}

func do_get_status(msg string){
	fmt.Println("do_get_status:", msg)
}

func do_get_paulse_status(msg string){
	fmt.Println("do_get_paulse_status:", msg)
}

func do_get_paulse_count(msg string){
	fmt.Println("do_get_paulse_count:", msg)
}

func GetRouter() *chi.Mux {
    router := chi.NewRouter()
    return router
}