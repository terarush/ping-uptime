package entity

type ChartDataPoint struct {
	Date    string  `json:"date"`
	Status  string  `json:"status"` // up, down, unknown
	Uptime  float64 `json:"uptime_pct"`
	Failed  int     `json:"failed"`
	Total   int     `json:"total"`
	Latency float64 `json:"latency"`
}

type MonitorStats struct {
	MonitorID    uint             `json:"monitor_id"`
	MonitorName  string           `json:"monitor_name"`
	MonitorURL   string           `json:"monitor_url"`
	Window       string           `json:"window"` // 1h, 1d, 1w, 1m, 1y, all
	UptimePct    float64          `json:"uptime_pct"`
	TotalChecks  int              `json:"total_checks"`
	FailedChecks int              `json:"failed_checks"`
	Points       []ChartDataPoint `json:"points"`
	Status       string           `json:"status"` // operational, degraded, outage
	AvgLatency   float64          `json:"avg_latency"`
	MinLatency   float64          `json:"min_latency"`
	MaxLatency   float64          `json:"max_latency"`
}
