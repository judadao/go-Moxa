package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := setupServer()
	http.ListenAndServe(":3000", r)
}

func setupServer() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homeHandler) // 接收IP

	r.Mount("/moxa", moxaRoutes())

	fs := http.FileServer(http.Dir("static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 導向到 HTML 表單頁面
	http.Redirect(w, r, "/static/rest_view.html", http.StatusSeeOther)
}



func moxaRoutes() chi.Router {
	r := chi.NewRouter()
	diHandler := DiHandler{
	}
	// 靜態檔案服務
	r.Get("/getdi", diHandler.Get_di) // 傳送IP
	
	

	return r
}