package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

func ReportHandler(writer http.ResponseWriter, request *http.Request, mut *sync.Mutex) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем массив строк из параметра "strings"
	str := request.Form["strings"]
	args := strings.Join(str, " ")

	mut.Lock()

	conn, err := net.Dial("tcp", "localhost:9090")
	defer func(conn net.Conn) {
		_ = conn.Close()
		mut.Unlock()
	}(conn)

	_, err = fmt.Fprint(conn, "create_report "+args)
	if err != nil {
		return
	}

	// Получаем о обрабатываем результат.
	decoder := json.NewDecoder(conn)
	var receivedData map[string]interface{}
	if err := decoder.Decode(&receivedData); err != nil {
		log.Println("[✗] Error: there was an error decoding received data.")
		_, err := fmt.Fprintf(writer, "[✗] Error: the report could not be displayed, contact your system administrator.")
		if err != nil {
			return
		}
		return
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&receivedData); err != nil {
		log.Println("[✗] Error: there was an error sending data to the client.")
	}
}
