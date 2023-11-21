package handlers

import (
	"Shortener/dbms"
	"errors"
	"strings"
)

// DatabaseHandler QUERY (add <SHORT_URL> <URL>)
func DatabaseHandler(file string, query string) (string, error) {
	request := strings.Fields(query)

	switch strings.ToLower(request[0]) {
	case "add":
		if len(request) != 3 {
			return "", errors.New("invalid request")
		}

		result, err := dbms.AddToDatabase(file, request[1], request[2])
		if err != nil {
			return "", err
		}

		return result, nil

	case "get":
		if len(request) != 2 {
			return "", errors.New("invalid request")
		}

		result, err := dbms.FindInDatabase(file, request[1])
		if err != nil {
			return "", err
		}

		return result, nil

	default:
		return "", errors.New("invalid request")
	}
}
