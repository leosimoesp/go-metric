package repository

import (
	"sync"
	"testing"
	"time"

	"github.com/leosimoesp/go-metric/internal/app/metricdata"
)

func Test_storage_RemoveOld(t *testing.T) {

	db := map[string][]*metricdata.Metric{}

	now := time.Now()
	before := now.Add(-3601 * time.Second)

	m1 := &metricdata.Metric{ID: before.Unix(), Value: int64(20), CreatedAt: before}

	m2 := &metricdata.Metric{ID: now.Unix(), Value: int64(5), CreatedAt: now}

	db["active_visitors"] = []*metricdata.Metric{m1, m2}

	type fields struct {
		RWMutex sync.RWMutex
		values  map[string][]*metricdata.Metric
	}
	type args struct {
		key string
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "[1]-Remove only the record after one hour",
			fields: fields{
				values: db,
			},
			args: args{
				key: "active_visitors",
				id:  now.Unix(),
			},
			want: int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				RWMutex: tt.fields.RWMutex,
				values:  tt.fields.values,
			}
			got, err := s.RemoveOld(tt.args.key, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("storage.RemoveOld() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.RemoveOld() = %v, want %v", got, tt.want)
			}
		})
	}
}
