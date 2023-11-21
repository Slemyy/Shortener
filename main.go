package main

import (
	"Shortener/config"
	"Shortener/handlers"
	"log"
	"net"
	"net/http"
	"sync"
)

var conn net.Conn
var mut sync.Mutex

func main() {
	// Загружаем конфиг для работы с программой.
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	//// Загружаем базу данных для работы с программой.
	//conn, err = dbms.LoadDBMS(cfg)
	//if err != nil {
	//	log.Fatalln("Error connecting to server:", err.Error())
	//	return
	//}

	log.Println("[✔] Database loaded successfully.")

	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", redirect)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error loading database:", err.Error())
	}
}

func shorten(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handlers.ShortenHandler(writer, request, &mut)
}

func redirect(writer http.ResponseWriter, request *http.Request) {
	handlers.RedirectHandler(writer, request, conn, &mut)
}
