package handlers

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
)

func ShortenHandler(writer http.ResponseWriter, request *http.Request, mut *sync.Mutex) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	shortURL := generateShortURL(6)

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return
	}

	defer conn.Close()

	fmt.Fprint(conn, "--file database --query 'add "+shortURL+" "+originalURL+"'")
	req, _ := bufio.NewReader(conn).ReadString('\n')

	_, err = fmt.Fprintf(writer, "Shortened URL: http://localhost:8080/%s", req)
	if err != nil {
		return
	}
}

// Генерация случайного сокращенного URL
func generateShortURL(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
