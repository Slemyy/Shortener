package main

import (
	"Shortener/config"
	h "Shortener/handlers"
	"log"
	"net"
	"net/http"
)

type Service struct {
	cfg *config.Config
}

func main() {
	// Загружаем конфиг для работы с программой.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Загружаем базу данных для работы с программой.
	_, err = net.Dial(cfg.Database.Network, cfg.Database.Port)
	if err != nil {
		log.Fatalln("Error connecting to server:", err.Error())
		return
	}

	log.Println("[✔] Database loaded successfully.")

	http.HandleFunc("/", h.RedirectHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error loading database:", err.Error())
	}
}
