package database

import (
	"errors"
	"os"
	"strings"
)

const arraySize = 769

type HashNode struct {
	shortURL    string
	originalURL string
}

type HashTable struct {
	database [arraySize]*HashNode
}

func hashFunction(key string) int {
	hash := 0
	for i := 0; i < len(key); i++ {
		hash += int(key[i])
	}

	return hash % arraySize
}

func (hashMap *HashTable) Clear() error {
	for i := 0; i < arraySize; i++ {
		if hashMap.database[i] != nil {
			hashMap.database[i] = nil
		}
	}
	return nil
}

func (hashMap *HashTable) Insert(shortURL string, originalURL string) (string, error) {
	dataBaseEntry := &HashNode{shortURL, originalURL}
	index := hashFunction(shortURL)

	if hashMap.database[index] == nil {
		hashMap.database[index] = dataBaseEntry
		return hashMap.database[index].shortURL, nil
	}

	if hashMap.database[index].shortURL == shortURL {
		return hashMap.database[index].shortURL, nil
	}

	for i := (index + 1) % arraySize; i != index; i = (i + 1) % arraySize {
		if i == index {
			return "", errors.New("database is full")
		}

		if hashMap.database[i] == nil {
			hashMap.database[i] = dataBaseEntry
			return hashMap.database[index].shortURL, nil
		}

		if hashMap.database[i].shortURL == shortURL {
			return hashMap.database[index].shortURL, nil
		}
	}

	return "", errors.New("failed to add element")
}

func (hashMap *HashTable) Remove(shortURL string) (string, error) {
	index := hashFunction(shortURL)

	if hashMap.database[index] == nil {
		return "", errors.New("no such meaning")
	}

	// Если найден ключ в текущем индексе, удаляем его
	if hashMap.database[index].shortURL == shortURL {
		remove := hashMap.database[index].originalURL
		hashMap.database[index] = nil
		return remove, nil
	}

	// В случае коллизии, ищем ключ в следующих индексах
	for i := (index + 1) % arraySize; i != index; i = (i + 1) % arraySize {
		if hashMap.database[i] == nil {
			return "", errors.New("element not found")
		}

		if hashMap.database[i].shortURL == shortURL {
			remove := hashMap.database[i].originalURL
			hashMap.database[i] = nil
			return remove, nil
		}
	}

	return "", errors.New("no such meaning")
}

func (hashMap *HashTable) Get(shortURL string) (string, error) {
	index := hashFunction(shortURL)

	if hashMap.database[index] == nil {
		return "", errors.New("element not found")
	}

	if hashMap.database[index].shortURL == shortURL {
		return hashMap.database[index].originalURL, nil
	}

	for i := (index + 1) % arraySize; i != index; i = (i + 1) % arraySize {
		if hashMap.database[i] == nil {
			return "", errors.New("element not found")
		}

		if hashMap.database[i].shortURL == shortURL {
			return hashMap.database[i].originalURL, nil
		}
	}

	return "", errors.New("element not found")
}

// ReadFromFile Чтение дата файла.
func (hashMap *HashTable) ReadFromFile(filename string) error {
	content, err := os.ReadFile(filename)

	if err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			key := parts[0]
			value := strings.Join(parts[1:], " ")
			_, err = hashMap.Insert(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteToFile Запись в дата файл
func (hashMap *HashTable) WriteToFile(hashFile string) error {
	file, err := os.Create(hashFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	for i := 0; i < arraySize; i++ {
		if hashMap.database[i] != nil {
			_, err = file.WriteString(hashMap.database[i].shortURL + " " + hashMap.database[i].originalURL + "\n")
			if err != nil {
				return err
			}
		}
	}

	err = hashMap.Clear()
	if err != nil {
		return err
	}

	return nil
}
