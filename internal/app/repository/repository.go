package repository

import (
	"sync"
	"time"

	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/pkg/timehelper"
)

//go:generate mockery --name=MetricRepo --output=./mocks
type MetricRepo interface {
	Get(key string, id int64) (*metricdata.Metric, error)
	GetAllGreaterThan(key string, id int64) ([]*metricdata.Metric, error)
	RemoveOld(key string, id int64) (int64, error)
	Save(key string, metric metricdata.Metric) (int64, error)
}

type storage struct {
	sync.RWMutex
	values map[string][]*metricdata.Metric
}

func NewMetricRepo() MetricRepo {
	return storage{
		values: make(map[string][]*metricdata.Metric),
	}
}

func (s storage) Get(key string, id int64) (*metricdata.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	if records, ok := s.values[key]; ok && len(records) > 0 {
		for _, metric := range records {
			if metric.ID == id {
				return metric, nil
			}
		}
	}

	return nil, nil
}

func (s storage) GetAllGreaterThan(key string, id int64) ([]*metricdata.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	metrics := []*metricdata.Metric{}

	if records, ok := s.values[key]; ok && len(records) > 0 {
		for _, metric := range records {
			if metric.ID > id {
				metrics = append(metrics, metric)
			}
		}
	}

	return metrics, nil
}

func (s storage) RemoveOld(key string, id int64) (int64, error) {
	s.Lock()
	defer s.Unlock()

	if records, ok := s.values[key]; ok && len(records) > 0 {
		lastIndex := -1
		for k, metric := range records {
			if metric.ID < id {
				lastIndex = k
			}
		}
		if lastIndex > 0 && lastIndex+1 < len(records) {
			newValues := records[lastIndex+1:]
			s.values[key] = newValues
			return int64(lastIndex + 1), nil
		}
	}

	return int64(0), nil
}

func (s storage) Save(key string, metric metricdata.Metric) (int64, error) {
	s.Lock()
	defer s.Unlock()

	now := time.Now()

	id := timehelper.TimestampFromTime(now)
	metric.CreatedAt = now
	metric.ID = id

	if records, ok := s.values[key]; ok && len(records) > 0 {
		s.values[key] = append(s.values[key], &metric)
	} else {
		s.values[key] = []*metricdata.Metric{&metric}
	}

	return id, nil
}
