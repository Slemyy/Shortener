package main

import (
	"Shortener/handlers"
	"log"
	"net/http"
	"sync"
)

var mut sync.Mutex

func main() {
	log.Println("[✔] Database loaded successfully.")

	http.HandleFunc("/home", index)
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", redirect)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error loading database:", err.Error())
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	handlers.IndexHandler(writer, request)
}

func shorten(writer http.ResponseWriter, request *http.Request) {
	handlers.ShortenHandler(writer, request, &mut)
}

func redirect(writer http.ResponseWriter, request *http.Request) {
	handlers.RedirectHandler(writer, request, &mut)
}
