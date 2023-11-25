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

	_, err = fmt.Fprint(conn, "--file database/db --query 'get "+shortURL+"'")
	if err != nil {
		return
	}
	req, _ := bufio.NewReader(conn).ReadString('\n')

	if req[:5] != "Error" {
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
