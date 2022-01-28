package service

import (
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/leosimoesp/go-metric/config"
	"github.com/leosimoesp/go-metric/internal/app/constants"
	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/internal/app/repository"
	"github.com/leosimoesp/go-metric/pkg/errorwrapper"
	"github.com/leosimoesp/go-metric/pkg/log"
	"github.com/leosimoesp/go-metric/pkg/timehelper"
)

type MetricService interface {
	Save(key string, metricInput metricdata.MetricInput) (int64, error)
	CalculateSumMetrics(key string) (int64, error)
}

type metricInst struct {
	metricRepo repository.MetricRepo
	scheduler  *gocron.Scheduler
}

func NewMetricService() MetricService {

	repo := repository.NewMetricRepo()
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.StartAsync()

	return metricInst{
		metricRepo: repo,
		scheduler:  scheduler,
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

	scheduledTimeInMin, err := strconv.Atoi(os.Getenv(config.RemoveMetricsIntervalInMin))

	if err != nil {
		log.Logger().Errorf("CalculateSumMetrics %v", err)
		return 0, err
	}

	m.scheduler.Every(scheduledTimeInMin).Minute().Do(func() {
		m.metricRepo.RemoveOld(key, startID)
	})

	return sum, nil
}
