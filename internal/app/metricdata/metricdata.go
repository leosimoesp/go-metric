package metricdata

import "time"

type Metric struct {
	ID        int64
	Value     int64
	CreatedAt time.Time
}

type MetricInput struct {
	Value int `json:"value"`
}

type MetricSumResponse struct {
	Value int `json:"value"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
