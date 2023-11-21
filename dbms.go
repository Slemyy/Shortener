package main

import (
	"Shortener/config"
	"Shortener/handlers"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	// Загружаем конфиг для работы с программой.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	listener, err := net.Listen(cfg.Database.Network, cfg.Database.Port)
	if err != nil {
		log.Fatalln("Error starting database:", err.Error())
		return
	}

	defer func(listener net.Listener) { _ = listener.Close() }(listener)
	log.Println("The database was loaded on the port", listener.Addr().String()[5:]+"...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error connecting to database", err.Error())
			return
		}

		var mut sync.Mutex
		go handleClient(conn, &mut, cfg)
	}
}

func handleClient(conn net.Conn, mut *sync.Mutex, cfg *config.Config) {
	defer func(conn net.Conn) { _ = conn.Close() }(conn)
	remoteAddr := conn.RemoteAddr() // Получение адреса удаленного узла

	log.Printf("Connection established with: %s\n", remoteAddr)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Connection to %s is closed.\n", remoteAddr)
			return
		}

		clientMessage := string(buffer[:n])
		log.Printf("Received from %s: %s", remoteAddr, clientMessage)
		args := strings.Fields(clientMessage)

		if len(args) < 4 {
			_, err = conn.Write([]byte("Not enough arguments. Use: --file <file.json> --query <query>.\n"))
			if err != nil {
				log.Printf("(%s) Error: %v\n", remoteAddr, err)
				break
			}

			continue

		} else if args[0] != "--file" || args[2] != "--query" {
			_, err = conn.Write([]byte("Not enough arguments. Use: --file <file.json> --query <query>.\n"))
			if err != nil {
				log.Printf("(%s) Error: %v\n", remoteAddr, err)
				break
			}

			continue
		}

		query := strings.Join(args[3:], " ")

		if query[0] == '\'' || query[0] == '"' || query[0] == '<' {
			query = query[1:] // Убираем лишние элементы
		}

		if query[len(query)-1] == '\'' || query[len(query)-1] == '"' || query[len(query)-1] == '>' {
			query = query[:len(query)-1]
		}

		ans, err := handlers.DatabaseHandler(args[1], query, mut)
		if err != nil {
			response := "Error: " + err.Error() + "\n"
			_, err := conn.Write([]byte(response))
			if err != nil {
				log.Printf("(%s) Error: %v\n", remoteAddr, err)
				break
			}
		}

		// Отправка ответа клиенту
		log.Printf("[✔] (%s) Request processed successfully.", remoteAddr)
		_, err = conn.Write([]byte(ans + "\n"))
		if err != nil {
			log.Printf("(%s) Error: %v\n", remoteAddr, err)
			break
		}
	}
}
