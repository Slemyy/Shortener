package reporter

import (
	"encoding/json"
	"os"
)

func WriteJSONToFile(data interface{}, filename string) error {
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

func CreateReport(hierarchy []string, statistics []Statistic) map[string]interface{} {
	report := make(map[string]interface{})

	for _, stats := range statistics {
		if stats.PID == 0 {
			continue
		}

		URL := findURLByID(stats.PID, statistics)
		IP := stats.SourceIP
		TimeInterval := stats.Time[11:]

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
