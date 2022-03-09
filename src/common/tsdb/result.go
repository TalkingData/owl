package tsdb

type Result struct {
	Metric        string             `json:"metric"`
	Tags          map[string]string  `json:"tags"`
	AggregateTags []string           `json:"aggregateTags"`
	Dps           map[string]float64 `json:"dps"`
}
