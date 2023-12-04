package reporter

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func UpdateStatistics(oldURL string, shortURL string, IP string) {
	URL := oldURL + " (" + shortURL + ")"

	parentStats := Statistic{
		URL:   URL,
		Count: 1,
	}

	newStats := Statistic{
		SourceIP: IP,
		Time:     time.Now().Format("02-01-2006 15:04"),
		Count:    1,
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
