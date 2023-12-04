package reporter

const statsFilename = "stats.json"
const ReportFilename = "report.json"

type Statistic struct {
	ID       int    `json:"id"`
	PID      int    `json:"pid"`
	URL      string `json:"url"`
	SourceIP string `json:"sourceIP"`
	Time     string `json:"time"`
	Count    int    `json:"count"`
}
