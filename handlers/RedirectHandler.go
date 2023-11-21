package handlers

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

func RedirectHandler(writer http.ResponseWriter, request *http.Request, mut *sync.Mutex) {
	mut.Lock()
	defer mut.Unlock()

	shortURL := strings.TrimPrefix(request.URL.Path, "/")

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return
	}

	defer conn.Close()

	fmt.Fprint(conn, "--file database --query 'get "+shortURL+"'")
	req, _ := bufio.NewReader(conn).ReadString('\n')

	if req != "" {
		http.Redirect(writer, request, req[:len(req)-1], http.StatusFound)
	} else {
		http.NotFound(writer, request)
	}
}
