package handlers

import (
	"fmt"
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
	defer mut.Unlock()

	conn, err := net.Dial("tcp", "localhost:9090")
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	_, err = fmt.Fprint(conn, "create_report "+args)
	if err != nil {
		return
	}
}
