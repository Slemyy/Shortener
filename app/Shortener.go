package main

import (
	"Shortener/config"
	"Shortener/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var mut sync.Mutex

func main() {
	cfg, err := config.LoadConfig() // Загрузка конфига
	if err != nil {
		log.Fatalln("[✗] Error loading config:", err.Error())
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/{shortURL}", func(writer http.ResponseWriter, request *http.Request) {
		handlers.RedirectHandler(writer, request, &mut)
	}).Methods("GET")
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/shorten", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShortenHandler(writer, request, &mut)
	}).Methods("POST")

	http.HandleFunc("/styles.css", ServeCSS)
	http.Handle("/", r)

	log.Println("[✔] Server started successfully.")
	err = http.ListenAndServe(cfg.Shortener.Port, nil)
	if err != nil {
		log.Fatalln("[✗] Service loading error:", err.Error())
		return
	}

	log.Println("[✔] Server stopped.")
}

// ServeCSS подключение CSS стилей.
func ServeCSS(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "website/styles.css")
}

type PageVariables struct {
	Title string
}

func HomePage(writer http.ResponseWriter, _ *http.Request) {
	pageVariables := PageVariables{
		Title: "Shortener",
	}

	tmpl, err := template.ParseFiles("website/index.html")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(writer, pageVariables)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
