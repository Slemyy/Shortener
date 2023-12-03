package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const statsFilename = "stats.json"

type Statistic struct {
	ID           int    `json:"id"`
	PID          int    `json:"pid"`
	URL          string `json:"url"`
	SourceIP     string `json:"sourceIP"`
	TimeInterval string `json:"time"`
	Count        int    `json:"count"`
}

func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("[✗] Error starting database:", err.Error())
		return
	}

	defer func(listener net.Listener) { _ = listener.Close() }(listener)
	log.Println("[✔] The statistics service is running on the port", listener.Addr().String()[5:]+"...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("[✗] Error connecting to database", err.Error())
			return
		}

		var mutex sync.Mutex
		go handleConnection(conn, &mutex)
	}
}

func handleConnection(conn net.Conn, s *sync.Mutex) {
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		clientMessage := string(buffer[:n])
		log.Printf("Service request: %s", clientMessage)
		args := strings.Fields(clientMessage)

		if args[0] == "add_stats" { // Добавляем статистику
			addStatistics(args[1], args[2], args[3])
		} else if args[0] == "create_report" { // Формируем отчет
			var hierarchy []string

			for i := 1; i < len(args); i++ {
				hierarchy = append(hierarchy, args[i])
			}

			connDB, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Fatalln("[✗] Error connecting to database.")
				return
			}
			defer func(connDB net.Conn) { _ = conn.Close() }(connDB)

			_, err = connDB.Write([]byte("--file none --query \"create_report\"\n"))
			if err != nil {
				return
			}

			response, err := bufio.NewReader(connDB).ReadBytes(']')
			if err != nil {
				log.Fatalln("[✗] Error reading response from database server:", err)
				return
			}

			JsonFile := ByteToJson(response)
			JsonData := createReport(hierarchy, JsonFile)

			err = writeJSONToFile(JsonData, "report.json")
			if err != nil {
				fmt.Println("[✗] Error writing to report file:", err)
				return
			}

			log.Println("[✔] Report created successfully.")

		} else {
			_, err = conn.Write([]byte("[✗] Invalid request. Commands: {add_stats, create_report}\n"))
			log.Fatalln("[✗] Error: invalid request")
		}
	}
}

func writeJSONToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func createReport(hierarchy []string, statistics []Statistic) map[string]interface{} {
	report := make(map[string]interface{})

	for _, stats := range statistics {
		if stats.PID == 0 {
			continue
		}

		URL := findURLByID(stats.PID, statistics)
		IP := stats.SourceIP
		TimeInterval := stats.TimeInterval[11:]

		currLevel := report
		for _, level := range hierarchy {
			if level == "SourceIP" {
				if _, ok := currLevel[IP]; !ok {
					currLevel[IP] = make(map[string]interface{})
					if _, ok := currLevel["Sum"]; !ok {
						currLevel["Sum"] = 0
					}
				}
				currLevel = currLevel[IP].(map[string]interface{})
			} else if level == "TimeInterval" {
				if _, ok := currLevel[TimeInterval]; !ok {
					currLevel[TimeInterval] = make(map[string]interface{})
					if _, ok := currLevel["Sum"]; !ok {
						currLevel["Sum"] = 0
					}
				}
				currLevel = currLevel[TimeInterval].(map[string]interface{})
			} else if level == "URL" {
				if _, ok := currLevel[URL]; !ok {
					currLevel[URL] = make(map[string]interface{})
					if _, ok := currLevel["Sum"]; !ok {
						currLevel["Sum"] = 0
					}
				}
				currLevel = currLevel[URL].(map[string]interface{})
			}

			if _, ok := currLevel["Sum"]; !ok {
				currLevel["Sum"] = 0
			}
			currLevel["Sum"] = currLevel["Sum"].(int) + 1
		}

	}

	delete(report, "Sum")
	return report
}

func findURLByID(id int, statistics []Statistic) string {
	for _, stats := range statistics {
		if stats.ID == id {
			return stats.URL
		}
	}

	return ""
}

func ByteToJson(file []byte) []Statistic {
	var statistics []Statistic

	if len(file) == 0 {
		return nil
	}

	err := json.Unmarshal(file, &statistics)
	if err != nil {
		return nil
	}

	return statistics
}

func addStatistics(oldURL string, shortURL string, IP string) {
	URL := oldURL + " (" + shortURL + ")"

	parentStats := Statistic{
		URL:   URL,
		Count: 1,
	}

	newStats := Statistic{
		SourceIP:     IP,
		TimeInterval: time.Now().Format("02-01-2006 15:04"),
		Count:        1,
	}

	statistics, err := readStatisticsFromFile()
	if err != nil {
		log.Fatalln("[✗] Error reading statistics file.")
	}

	if statistics == nil {
		statistics = []Statistic{}
	}

	parentStats.ID = genUniqueID(statistics)
	if UniqueParents(statistics, parentStats.URL) == true {
		statistics = append(statistics, parentStats)
	} else {
		ParentsCount(statistics, parentStats.URL)
	}

	newStats.ID = genUniqueID(statistics)
	newStats.PID = genPID(statistics, URL)
	statistics = append(statistics, newStats)

	err = writeStatisticsToFile(statistics)
}

func writeStatisticsToFile(statistics []Statistic) error {
	jsonData, err := json.MarshalIndent(statistics, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(statsFilename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func genPID(statistics []Statistic, url string) int {
	PID := 0
	for _, stats := range statistics {
		if stats.URL == url {
			PID = stats.ID
		}
	}

	return PID
}

func ParentsCount(statistics []Statistic, url string) {
	for index := range statistics {
		if statistics[index].URL == url {
			statistics[index].Count++
			return
		}
	}
}

func UniqueParents(statistics []Statistic, url string) bool {
	for _, stats := range statistics {
		if stats.URL == url {
			return false
		}
	}

	return true
}

func genUniqueID(statistics []Statistic) int {
	maxID := 0

	for _, stats := range statistics {
		if stats.ID > maxID {
			maxID = stats.ID
		}
	}

	return maxID + 1
}

func readStatisticsFromFile() ([]Statistic, error) {
	var statistics []Statistic

	file, err := os.ReadFile(statsFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(file, &statistics)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}
