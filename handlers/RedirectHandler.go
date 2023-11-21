package handlers

import (
	"net"
	"net/http"
	"sync"
)

func RedirectHandler(writer http.ResponseWriter, request *http.Request, conn net.Conn, mut *sync.Mutex) {

}
