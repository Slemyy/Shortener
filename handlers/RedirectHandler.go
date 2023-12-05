package handlers

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"strings"
	"sync"
)

func RedirectHandler(writer http.ResponseWriter, request *http.Request, mut *sync.Mutex) {
	mut.Lock()

	vars := mux.Vars(request)
	shortURL := vars["shortURL"]

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
		mut.Unlock()
	}(conn)

	_, err = fmt.Fprint(conn, "--file database/db.data --query 'get "+shortURL+"'")
	if err != nil {
		return
	}
	req, _ := bufio.NewReader(conn).ReadString('\n')

	if req[:5] != "Error" {
		go func() {
			connReport, _ := net.Dial("tcp", "localhost:9090")
			_, err = fmt.Fprint(connReport, "add_stats "+req[:len(req)-1]+" "+shortURL+" "+GetClientIP(request)+"\n")
			_ = connReport.Close()
		}()

		switch strings.HasPrefix(req, "http") || strings.HasPrefix(req, "https") {
		case true:
			http.Redirect(writer, request, req[:len(req)-1], http.StatusFound)
		case false:
			http.Redirect(writer, request, "http://"+req[:len(req)-1], http.StatusFound)
		}
	} else {
		http.Redirect(writer, request, "http://localhost:8080/", http.StatusFound)
	}
}

func GetClientIP(request *http.Request) string {
	// Получаем IP-адрес клиента из заголовка X-Real-IP или X-Forwarded-For
	ip := request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = request.Header.Get("X-Forwarded-For")
	}

	// Если заголовки не содержат информацию о IP, используем RemoteAddr
	if ip == "" {
		ip, _, _ = net.SplitHostPort(request.RemoteAddr)
	}

	return ip
}
