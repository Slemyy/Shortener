package dbms

import (
	"Shortener/config"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
	"time"
)

type Database struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	URL     string    `json:"URL"`
	Visits  int       `json:"visited"`
	Created time.Time `json:"created"`
}

func LoadDBMS(cfg *config.Config) (net.Conn, error) {
	conn, err := net.Dial(cfg.Database.Network, cfg.Database.Port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func hasSuffix(dataFile string) string {
	if !strings.HasSuffix(dataFile, ".json") {
		dataFile += ".json"
	}

	return dataFile
}

func readFromDatabase(dataFile string) ([]Database, error) {
	dataContent, err := os.ReadFile(dataFile)
	var shortURLs []Database

	if err == nil {
		err = json.Unmarshal(dataContent, &shortURLs)
		if err != nil {
			return nil, errors.New("could not read databases")
		}
	}

	return shortURLs, nil
}

func FindInDatabase(dataFile string, shortURL string, URL string) (string, error) {
	dataFile = hasSuffix(dataFile)

	DB, err := readFromDatabase(dataFile)
	if err != nil {
		return "", err
	}

	newShortURL := Database{
		Id:      len(DB),
		Name:    shortURL,
		URL:     URL,
		Visits:  0,
		Created: time.Now(),
	}

	// Проверка на уникальность новых данных в существующем срезе.
	isUnique := true
	for _, existingShortURL := range DB {
		if newShortURL.URL == existingShortURL.URL {
			isUnique = false
			// Если данные не уникальны, вернем значение.
			return existingShortURL.Name, nil
		}
	}

	// Если данные уникальны, добавьте их в существующий срез.
	if isUnique {
		DB = append(DB, newShortURL)
	}

	// Сериализация обновленных данных в JSON.
	jsonData, err := json.MarshalIndent(DB, "", "  ") // Добавляем отступы с двумя пробелами
	if err != nil {
		return "", err
	}

	// Сохранение JSON-данных в файл, перезаписывая существующий файл.
	err = os.WriteFile(dataFile, jsonData, 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}
