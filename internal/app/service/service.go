package service

import (
	"sync/atomic"

	"github.com/leosimoesp/go-metric/internal/app/constants"
	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/internal/app/repository"
	"github.com/leosimoesp/go-metric/pkg/errorwrapper"
	"github.com/leosimoesp/go-metric/pkg/timehelper"
)

type MetricService interface {
	Save(key string, metricInput metricdata.MetricInput) (int64, error)
	CalculateSumMetrics(key string) (int64, error)
}

type metricInst struct {
	metricRepo repository.MetricRepo
}

func NewMetricService() MetricService {
	return metricInst{
		metricRepo: repository.NewMetricRepo(),
	}
}

func (m metricInst) Save(key string, metricInput metricdata.MetricInput) (int64, error) {

	newID, err := m.metricRepo.Save(key, metricdata.Metric{
		Value: int64(metricInput.Value),
	})

	if err != nil {
		return 0, errorwrapper.NewErrorWrapper(constants.GenericDatabaseError,
			constants.GenericStatusCodeDatabaseError, err)
	}

	return newID, nil
}

func (m metricInst) CalculateSumMetrics(key string) (int64, error) {

	startID := timehelper.GetLastHourTimestamp()

	metrics, err := m.metricRepo.GetAllGreaterThan(key, startID)

	if err != nil {
		return 0, errorwrapper.NewErrorWrapper(constants.GenericDatabaseError,
			constants.GenericStatusCodeDatabaseError, err)
	}

	var sum int64

	for _, metric := range metrics {
		atomic.AddInt64(&sum, metric.Value)
	}

	return sum, nil
}
