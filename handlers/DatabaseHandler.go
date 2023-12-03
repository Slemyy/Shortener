package handlers

import (
	"Shortener/database"
	"errors"
	"strings"
	"sync"
)

// DatabaseHandler QUERY (add <SHORT_URL> <URL>)
func DatabaseHandler(file string, query string, mut *sync.Mutex) (string, error) {
	request := strings.Fields(query)

	mut.Lock()
	defer mut.Unlock()

	hash := &database.HashTable{}

	switch strings.ToLower(request[0]) {
	case "add":
		if len(request) != 3 {
			return "", errors.New("invalid request")
		}

		err := hash.ReadFromFile(file)
		if err != nil {
			return "", err
		}

		result, err := hash.Insert(request[1], request[2])
		if err != nil {
			return "", err
		}

		err = hash.WriteToFile(file)
		if err != nil {
			return "", err
		}

		return result, nil

	case "get":
		if len(request) != 2 {
			return "", errors.New("invalid request")
		}

		err := hash.ReadFromFile(file)
		if err != nil {
			return "", err
		}

		result, err := hash.Get(request[1])
		if err != nil {
			return "", err
		}

		err = hash.WriteToFile(file)
		if err != nil {
			return "", err
		}

		return result, nil

	default:
		return "", errors.New("invalid request")
	}
}
