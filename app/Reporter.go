package main

import (
	"Shortener/reporter"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("[✗] Error starting database:", err.Error())
		return
	}

	defer func(listener net.Listener) { _ = listener.Close() }(listener)
	log.Println("[✔] The statistics service is running on the port", listener.Addr().String()[5:]+"...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("[✗] Error connecting to database", err.Error())
			return
		}

		var mutex sync.Mutex
		go handleConnection(conn, &mutex)
	}
}

func handleConnection(conn net.Conn, mut *sync.Mutex) {
	defer func(conn net.Conn) {
		_ = conn.Close()
		mut.Unlock()
	}(conn)

	mut.Lock()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		clientMessage := string(buffer[:n])
		log.Printf("Service request: %s", clientMessage)
		args := strings.Fields(clientMessage)

		if args[0] == "add_stats" { // Добавляем статистику
			reporter.UpdateStatistics(args[1], args[2], args[3])
			log.Println("[✔] Statistics updated successfully.")
		} else if args[0] == "create_report" { // Формируем отчет
			var hierarchy []string

			for i := 1; i < len(args); i++ {
				hierarchy = append(hierarchy, args[i])
			}

			connDB, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Fatalln("[✗] Error connecting to database.")
				return
			}

			_, err = connDB.Write([]byte("--file none --query \"create_report\"\n"))
			if err != nil {
				return
			}

			response, err := bufio.NewReader(connDB).ReadBytes(']')
			if err != nil {
				log.Fatalln("[✗] Error reading response from database server:", err)
				return
			}

			_ = connDB.Close() // Закрываем соединение

			JsonFile := reporter.ByteToJson(response)
			JsonData := reporter.CreateReport(hierarchy, JsonFile)

			go func() {
				// Отправка результата клиенту
				encoder := json.NewEncoder(conn)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(&JsonData); err != nil {
					log.Println("[✗] Error: there was an error sending data to the client.")
				}
			}()

			err = reporter.WriteJSONToFile(&JsonData, reporter.ReportFilename)
			if err != nil {
				fmt.Println("[✗] Error writing to report file:", err)
				return
			}

			log.Println("[✔] Report created successfully.")

		} else {
			_, err = conn.Write([]byte("[✗] Invalid request. Commands: {add_stats, create_report}\n"))
			log.Fatalln("[✗] Error: invalid request")
		}
	}
}
