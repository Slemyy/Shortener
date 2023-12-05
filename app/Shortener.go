package main

import (
	"Shortener/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var mut sync.Mutex

func main() {
	r := mux.NewRouter()

	r.Use(checkSecretPath)

	r.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "Что ты тут делаешь?")
		if err != nil {
			return
		}
	})
	r.HandleFunc("/{shortURL}", func(writer http.ResponseWriter, request *http.Request) {
		handlers.RedirectHandler(writer, request, &mut)
	}).Methods("GET")
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/shorten", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShortenHandler(writer, request, &mut)
	}).Methods("POST")
	r.HandleFunc("/report", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ReportHandler(writer, request, &mut)
	}).Methods("POST")

	http.HandleFunc("/styles.css", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "website/styles.css")
	})
	http.Handle("/", r)

	log.Println("[✔] Server started successfully.")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalln("[✗] Service loading error:", err.Error())
		return
	}

	log.Println("[✔] Server stopped.")
}

func checkSecretPath(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/favicon.ico" {
			http.Error(writer, "Access to <"+request.URL.Path+"> is forbidden.", http.StatusForbidden)
			return
		}

		handler.ServeHTTP(writer, request)
	})
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
