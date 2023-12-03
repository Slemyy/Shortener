package handlers

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"sync"
)

func ShortenHandler(writer http.ResponseWriter, request *http.Request, mut *sync.Mutex) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	originalURL := request.Form.Get("url")
	if originalURL == "" {
		http.Error(writer, "URL is required", http.StatusBadRequest)
		return
	}

	mut.Lock()
	defer mut.Unlock()

	shortURL := generateShortURL(originalURL) // Сокращаем ссылку

	conn, err := net.Dial("tcp", "localhost:6379")
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	_, err = fmt.Fprint(conn, "--file database/db.data --query 'add "+shortURL+" "+originalURL+"'")
	if err != nil {
		return
	}
	req, _ := bufio.NewReader(conn).ReadString('\n')

	if req[:5] != "Error" {
		_, err = fmt.Fprint(writer, "http://localhost:8080/"+req)
		if err != nil {
			return
		}
	} else {
		_, err = fmt.Fprint(writer, req)
		if err != nil {
			return
		}
	}
}

func generateShortURL(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	shortLink := hex.EncodeToString(hash.Sum(nil))

	return shortLink[:7]
}
