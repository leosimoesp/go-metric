// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	metricdata "github.com/leosimoesp/go-metric/internal/app/metricdata"
	mock "github.com/stretchr/testify/mock"
)

// MetricRepo is an autogenerated mock type for the MetricRepo type
type MetricRepo struct {
	mock.Mock
}

// Get provides a mock function with given fields: key, id
func (_m *MetricRepo) Get(key string, id int64) (*metricdata.Metric, error) {
	ret := _m.Called(key, id)

	var r0 *metricdata.Metric
	if rf, ok := ret.Get(0).(func(string, int64) *metricdata.Metric); ok {
		r0 = rf(key, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metricdata.Metric)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(key, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllGreaterThan provides a mock function with given fields: key, id
func (_m *MetricRepo) GetAllGreaterThan(key string, id int64) ([]*metricdata.Metric, error) {
	ret := _m.Called(key, id)

	var r0 []*metricdata.Metric
	if rf, ok := ret.Get(0).(func(string, int64) []*metricdata.Metric); ok {
		r0 = rf(key, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*metricdata.Metric)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(key, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveOld provides a mock function with given fields: key, id
func (_m *MetricRepo) RemoveOld(key string, id int64) (int64, error) {
	ret := _m.Called(key, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int64) int64); ok {
		r0 = rf(key, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(key, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: key, metric
func (_m *MetricRepo) Save(key string, metric metricdata.Metric) (int64, error) {
	ret := _m.Called(key, metric)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, metricdata.Metric) int64); ok {
		r0 = rf(key, metric)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, metricdata.Metric) error); ok {
		r1 = rf(key, metric)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}