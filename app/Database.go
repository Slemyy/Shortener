package main

import (
	"Shortener/handlers"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalln("[✗] Error starting database:", err.Error())
		return
	}

	defer func(listener net.Listener) { _ = listener.Close() }(listener)
	log.Println("[✔] The database was loaded on the port", listener.Addr().String()[5:]+"...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("[✗] Error connecting to database", err.Error())
			return
		}

		var mutex sync.Mutex
		go handleClient(conn, &mutex)
	}
}

func handleClient(conn net.Conn, mutex *sync.Mutex) {
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		clientMessage := string(buffer[:n])
		log.Printf("Service request: %s", clientMessage)
		args := strings.Fields(clientMessage)

		if len(args) < 4 {
			_, err = conn.Write([]byte("[✗] Not enough arguments. Use: --file <file.json> --query <query>.\n"))
			if err != nil {
				log.Printf("[✗] Error: %v\n", err)
				break
			}

			continue

		} else if args[0] != "--file" || args[2] != "--query" {
			_, err = conn.Write([]byte("[✗] Not enough arguments. Use: --file <file.json> --query <query>.\n"))
			if err != nil {
				log.Printf("[✗] Error: %v\n", err)
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

		if len(query) > 12 && query[:13] == "create_report" {
			jsonFile, err := os.ReadFile("stats.json")
			if err != nil {
				log.Printf("[✗] Error: %v\n", err)
				break
			}

			log.Printf("[✔] Request processed successfully.")
			_, err = conn.Write(jsonFile)
			if err != nil {
				log.Printf("[✗] Error: %v\n", err)
				break
			}

			return
		}

		ans, err := handlers.DatabaseHandler(args[1], query, mutex)
		if err != nil {
			response := "Error: " + err.Error() + "\n"
			_, err := conn.Write([]byte(response))
			if err != nil {
				log.Printf("[✗] Error: %v\n", err)
				break
			}
		}

		// Отправка ответа клиенту
		log.Printf("[✔] Request processed successfully.")
		_, err = conn.Write([]byte(ans + "\n"))
		if err != nil {
			log.Printf("[✗] Error: %v\n", err)
			break
		}
	}
}
