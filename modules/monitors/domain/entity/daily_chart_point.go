package entity

type DailyChartPoint struct {
	Date       string  `json:"date"`
	UptimePct  float64 `json:"uptime_pct"`
	Status     string  `json:"status"`
	Failed     int     `json:"failed"`
	Total      int     `json:"total"`
}
