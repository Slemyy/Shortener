package handlers

import (
	"Shortener/dbms"
	"errors"
	"strings"
	"sync"
)

// DatabaseHandler QUERY (add <SHORT_URL> <URL>)
func DatabaseHandler(file string, query string, mut *sync.Mutex) (string, error) {
	request := strings.Fields(query)

	mut.Lock()
	defer mut.Unlock()

	switch strings.ToLower(request[0]) {
	case "add":
		if len(request) != 3 {
			return "", errors.New("invalid request")
		}

		result, err := dbms.FindInDatabase(file, request[1], request[2])
		if err != nil {
			return "", err
		}

		return result, nil

	case "get":
		return "", nil

	default:
		return "", errors.New("invalid request")
	}
}
