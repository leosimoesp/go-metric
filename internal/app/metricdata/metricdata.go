package metricdata

import "time"

type Metric struct {
	ID        int64
	Value     int64
	CreatedAt time.Time
}
