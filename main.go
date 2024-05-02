package main

import (
	"fmt"
	"go-Moxa/di"
	"go-Moxa/do"
	"net/http"

	"github.com/go-chi/chi"
)

func choose_api(callback func() ){
	callback()
}
var router *chi.Mux
func main() {
	router = chi.NewRouter()
	router.Get("/", do.Do_choose_api)
	di.Di_choose_api(20)
	// choose_api(do.Do_choose_api)
	fmt.Println("Server is running on :8080")
    http.ListenAndServe(":8080", router)

}
// func db_create(){
// 	db, err := sql.Open("sqlite3", "./di_data.db")
// 	if err != nil {
// 		fmt.Println("無法打開數據庫:", err)
// 		return
// 	}
// 	defer db.Close()

// 	// 創建數據庫表
// 	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS di_data (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		slot INTEGER,
// 		di_index INTEGER,
// 		di_mode INTEGER,
// 		di_status INTEGER
// 	)`)
// 	if err != nil {
// 		fmt.Println("無法創建表:", err)
// 		return
// 	}
// }

// func setupServer() chi.Router {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	r.Get("/", homeHandler) // 接收IP

// 	r.Mount("/moxa", moxaRoutes())

// 	fs := http.FileServer(http.Dir("static"))
//     r.Handle("/static/*", http.StripPrefix("/static/", fs))
	

// 	return r
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	// 導向到 HTML 表單頁面
// 	http.Redirect(w, r, "/static/rest_view.html", http.StatusSeeOther)
// }



// func moxaRoutes() chi.Router {
// 	r := chi.NewRouter()
// 	diHandler := DiHandler{
// 	}
// 	// DbHandler := DbHandler{}
// 	// 靜態檔案服務
// 	r.Get("/getdi", diHandler.Get_di) // 傳送IP
// 	// r.Get("/viewdb", DbHandler.ViewDBHandler) // 傳送IP
	

	
	
	

// 	return r
// }