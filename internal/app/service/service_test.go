package service

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/leosimoesp/go-metric/config"
	"github.com/leosimoesp/go-metric/internal/app/constants"
	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/internal/app/repository"
	"github.com/leosimoesp/go-metric/internal/app/repository/mocks"
	"github.com/leosimoesp/go-metric/pkg/errorwrapper"
	"github.com/leosimoesp/go-metric/pkg/timehelper"
	"github.com/stretchr/testify/mock"
)

func init() {
	os.Setenv(config.RemoveMetricsIntervalInMin, "10")
}

func Test_metricInst_Save(t *testing.T) {

	now := time.Now()
	expectedID := timehelper.TimestampFromTime(now)

	type fields struct {
		metricRepo repository.MetricRepo
		scheduler  *gocron.Scheduler
	}
	type args struct {
		key         string
		metricInput metricdata.MetricInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "[1]-Should save a new metric with success into a empty repo",
			fields: fields{
				metricRepo: &mocks.MetricRepo{},
				scheduler:  gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
				metricInput: metricdata.MetricInput{
					Value: 30,
				},
			},
			want: expectedID,
		},
		{
			name: "[2]-Should save a new metric with success into a old no empty repo",
			fields: fields{
				metricRepo: &mocks.MetricRepo{},
				scheduler:  gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
				metricInput: metricdata.MetricInput{
					Value: 4,
				},
			},
			want: expectedID,
		},
		{
			name: "[3]-Should return err when database is down",
			fields: fields{
				metricRepo: &mocks.MetricRepo{},
				scheduler:  gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
				metricInput: metricdata.MetricInput{
					Value: 4,
				},
			},
			wantErr: errors.New("erro"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := metricInst{
				metricRepo: tt.fields.metricRepo,
				scheduler:  tt.fields.scheduler,
			}

			switch tt.name {

			case "[1]-Should save a new metric with success into a empty repo":
				m.metricRepo.(*mocks.MetricRepo).On("Save", mock.Anything, mock.Anything).Return(tt.want, nil)

			case "[2]-Should save a new metric with success into a old no empty repo":
				m.metricRepo.(*mocks.MetricRepo).On("Save", mock.Anything, mock.Anything).Return(tt.want, nil)
			case "[3]-Should return err when database is down":
				m.metricRepo.(*mocks.MetricRepo).On("Save", mock.Anything, mock.Anything).Return(int64(0), tt.wantErr)
			}

			got, err := m.Save(tt.args.key, tt.args.metricInput)

			if err != nil && tt.wantErr == nil || err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("metricInst.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("metricInst.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_print(t *testing.T) {
	nowTimestamp := time.Now().Unix()
	lastHourTimestamp := timehelper.GetLastHourTimestamp()

	fiveSecondsBefore := time.Now().Add(-5 * time.Second).Unix()

	fmt.Println(nowTimestamp, lastHourTimestamp, nowTimestamp-lastHourTimestamp, fiveSecondsBefore)
}

func Test_metricInst_CalculateSumMetrics(t *testing.T) {

	fiveSecondsBefore := time.Now().Add(-5 * time.Second).Unix()
	oneMinuteBefore := time.Now().Add(-1 * time.Minute).Unix()
	oneHourBefore := time.Now().Add(-3600 * time.Second).Unix()

	type fields struct {
		metricRepo repository.MetricRepo
		scheduler  *gocron.Scheduler
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "[1]-Should return zero if there isn't any metric at database",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			want: int64(0),
		},
		{
			name: "[2]-Should return error when database is down",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			wantErr: errorwrapper.NewErrorWrapper(constants.GenericDatabaseError,
				constants.GenericStatusCodeDatabaseError, errors.New("timeout")),
		},
		{
			name: "[3]-Should return a correct sum for one metric at database",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			want: int64(10),
		},
		{
			name: "[4]-Should return a correct sum for more than one metric at database",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			want: int64(50),
		},
		{
			name: "[5]-Should not return sum if metric is older than one hour",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			want: int64(10),
		},
		{
			name: "[6]-Error when try get scheduler time env",
			fields: fields{
				&mocks.MetricRepo{},
				gocron.NewScheduler(time.UTC),
			},
			args: args{
				key: "active_visitors",
			},
			wantErr: errors.New("erro"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := metricInst{
				metricRepo: tt.fields.metricRepo,
				scheduler:  tt.fields.scheduler,
			}

			switch tt.name {
			case "[1]-Should return zero if there isn't any metric at database":
				os.Setenv(config.RemoveMetricsIntervalInMin, "10")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{}, nil)

			case "[2]-Should return error when database is down":
				os.Setenv(config.RemoveMetricsIntervalInMin, "10")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return(nil, tt.wantErr)

			case "[3]-Should return a correct sum for one metric at database":
				os.Setenv(config.RemoveMetricsIntervalInMin, "10")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{{ID: fiveSecondsBefore,
					Value: int64(10)}}, nil)
			case "[4]-Should return a correct sum for more than one metric at database":
				os.Setenv(config.RemoveMetricsIntervalInMin, "10")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{{ID: oneMinuteBefore, Value: int64(10)},
					{ID: fiveSecondsBefore, Value: int64(40)}}, nil)
			case "[5]-Should not return sum if metric is older than one hour":
				os.Setenv(config.RemoveMetricsIntervalInMin, "10")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{{ID: oneHourBefore,
					Value: int64(10)}}, nil)
			case "[6]-Error when try get scheduler time env":
				os.Setenv(config.RemoveMetricsIntervalInMin, "a")
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return(nil, tt.wantErr)
			}

			got, err := m.CalculateSumMetrics(tt.args.key)
			if err != nil && tt.wantErr == nil || err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("metricInst.CalculateSumMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("metricInst.CalculateSumMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
