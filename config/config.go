package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Shortener ShortenerConfig `json:"shortener"`
	Database  DatabaseConfig  `json:"database"`
	Report    ReportConfig    `json:"report"`
}

// ShortenerConfig структура для параметров сервиса
type ShortenerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DatabaseConfig структура для параметров базы данных
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Network  string `json:"network"`
	Location string `json:"location"`
	Name     string `json:"name"`
}

// ReportConfig структура для параметров сервиса отчетов
type ReportConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadConfig() (*Config, error) {
	configData, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	log.Println("[✔] Config loaded successfully.")

	return &config, nil
}
