package repositories

import (
	"alerting-service/internal/usecases"
	"reflect"
	"testing"
)

func TestNewStorageRepository(t *testing.T) {
	tests := []struct {
		name string
		want usecases.MemStorage
	}{
		{
			name: "new repository test",
			want: &MemStorageImp{gauges: map[string]float64{}, counters: map[string]int{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewStorageRepository()

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}

func TestGetCounterMetric(t *testing.T) {
	storage := NewStorageRepository()
	storage.UpdateCounterMetric("temp", 25)

	tests := []struct {
		name, metricName string
		want             int
		wantOk           bool
	}{
		{
			name:       "valid get counter test",
			metricName: "temp",
			want:       25,
			wantOk:     true,
		},
		{
			name:       "invalid get counter test",
			metricName: "count",
			want:       0,
			wantOk:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want, ok := storage.GetCounterMetric(test.metricName)

			if ok != test.wantOk {
				t.Errorf("want ok: %v, got: %v", test.wantOk, ok)
			}

			if want != test.want {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}
func TestGetGaugeMetric(t *testing.T) {
	storage := NewStorageRepository()
	storage.UpdateGaugeMetric("temp", 25.2)

	tests := []struct {
		name, metricName string
		want             float64
		wantOk           bool
	}{
		{
			name:       "valid get gauge test",
			metricName: "temp",
			want:       25.2,
			wantOk:     true,
		},
		{
			name:       "invalid get gauge test",
			metricName: "count",
			want:       0,
			wantOk:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want, ok := storage.GetGaugeMetric(test.metricName)

			if ok != test.wantOk {
				t.Errorf("want ok: %v, got: %v", test.wantOk, ok)
			}

			if want != test.want {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}

func TestUpdateGaugeMetric(t *testing.T) {
	storage := NewStorageRepository()

	tests := []struct {
		name, metricName  string
		want, metricValue float64
	}{
		{
			name:        "valid update gauge test",
			metricName:  "temp",
			metricValue: 25.2,
			want:        25.2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.UpdateGaugeMetric(test.metricName, test.metricValue)
			want, _ := storage.GetGaugeMetric(test.metricName)

			if want != test.want {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	tests := []struct {
		name, metricName                string
		want, metricValue, initialValue int
	}{
		{
			name:         "valid update gauge first updating test",
			metricName:   "temp",
			initialValue: 0,
			metricValue:  25,
			want:         25,
		},
		{
			name:         "valid update gauge next updating test",
			metricName:   "temp",
			initialValue: 25,
			metricValue:  25,
			want:         50,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewStorageRepository()
			storage.UpdateCounterMetric(test.metricName, test.initialValue)

			storage.UpdateCounterMetric(test.metricName, test.metricValue)
			want, _ := storage.GetCounterMetric(test.metricName)

			if want != test.want {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}
