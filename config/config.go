package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	Shortlink ShortlinkConfig `json:"shortlink"`
}

// ServerConfig структура для параметров сервера
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DatabaseConfig структура для параметров базы данных
type DatabaseConfig struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Network  string `json:"network"`
}

// ShortlinkConfig структура для параметров сокращения ссылок
type ShortlinkConfig struct {
	Length     int    `json:"length"`
	Characters string `json:"characters"`
}

func LoadConfig() (*Config, error) {
	configData, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	log.Println("[✔] Config loaded successfully.")

	return &config, nil
}
