package timehelper

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTimestamp(t *testing.T) {

	expected := time.Now().Unix()

	tests := []struct {
		name string
		want int64
	}{
		{
			name: "[1]-Should create a timestamp value from time now",
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := CreateTimestamp(); got != tt.want {
				t.Errorf("CreateTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTimestampOld(t *testing.T) {

	now := CreateTimestamp()

	lastHour := GetLastHourTimestamp()

	d1 := CreateTimestamp()
	time.Sleep(1 * time.Second)
	d2 := CreateTimestamp()

	fmt.Println(d1 == d2, now, lastHour, now-lastHour)
	assert.Greater(t, d2, d1, "The d2 time should be greather than d1.")
}

func TestGetLastHourTimestamp(t *testing.T) {

	now := time.Now()

	lastHour := now.Add(-1 * time.Hour)

	expected := lastHour.Unix()

	tests := []struct {
		name string
		want int64
	}{
		{
			name: "[1]-Should get last one hour timestamp",
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLastHourTimestamp(); got != tt.want {
				t.Errorf("GetLastHourTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
