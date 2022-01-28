package service

import (
	"errors"
	"testing"
	"time"

	"github.com/leosimoesp/go-metric/internal/app/constants"
	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/internal/app/repository"
	"github.com/leosimoesp/go-metric/internal/app/repository/mocks"
	"github.com/leosimoesp/go-metric/pkg/errorwrapper"
	"github.com/leosimoesp/go-metric/pkg/timehelper"
	"github.com/stretchr/testify/mock"
)

func Test_metricInst_Save(t *testing.T) {

	now := time.Now()
	expectedID := timehelper.TimestampFromTime(now)

	type fields struct {
		metricRepo repository.MetricRepo
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

func Test_metricInst_CalculateSumMetrics(t *testing.T) {
	type fields struct {
		metricRepo repository.MetricRepo
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
			},
			args: args{
				key: "active_visitors",
			},
			want: int64(50),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := metricInst{
				metricRepo: tt.fields.metricRepo,
			}

			switch tt.name {
			case "[1]-Should return zero if there isn't any metric at database":
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{}, nil)

			case "[2]-Should return error when database is down":
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return(nil, tt.wantErr)

			case "[3]-Should return a correct sum for one metric at database":
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{{ID: int64(1), Value: int64(10)}}, nil)
			case "[4]-Should return a correct sum for more than one metric at database":
				m.metricRepo.(*mocks.MetricRepo).On("GetAllGreaterThan",
					mock.Anything, mock.Anything).Return([]*metricdata.Metric{{ID: int64(1), Value: int64(10)},
					{ID: int64(2), Value: int64(40)}}, nil)
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
