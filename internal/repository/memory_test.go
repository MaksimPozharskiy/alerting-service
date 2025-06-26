package repository

import (
	"alerting-service/internal/models"
	"alerting-service/internal/utils"
	"reflect"
	"testing"
)

func TestNewStorageRepository(t *testing.T) {
	tests := []struct {
		name string
		want StorageRepository
	}{
		{
			name: "new repository test",
			want: &MemStorageImp{gauges: map[string]float64{}, counters: map[string]int{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewMemStorageRepository()

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}

func TestGetCounterMetric(t *testing.T) {
	storage := NewMemStorageRepository()
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
			want, ok, _ := storage.GetCounterMetric(test.metricName)

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
	storage := NewMemStorageRepository()
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
			want, ok, _ := storage.GetGaugeMetric(test.metricName)

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
	storage := NewMemStorageRepository()

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
			want, _, _ := storage.GetGaugeMetric(test.metricName)

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
			storage := NewMemStorageRepository()
			storage.UpdateCounterMetric(test.metricName, test.initialValue)

			storage.UpdateCounterMetric(test.metricName, test.metricValue)
			want, _, _ := storage.GetCounterMetric(test.metricName)

			if want != test.want {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	storage := NewMemStorageRepository()
	_ = storage.UpdateGaugeMetric("gauge1", 1.23)
	_ = storage.UpdateCounterMetric("counter1", 10)

	metrics, err := storage.GetMetrics()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
}

func TestSetMetrics(t *testing.T) {
	storage := NewMemStorageRepository()

	metrics := []models.Metrics{
		{ID: "g1", MType: models.GaugeMetric, Value: utils.FloatPtr(5.5)},
		{ID: "c1", MType: models.CounterMetric, Delta: utils.IntPtr(15)},
	}

	storage.SetMetrics(metrics)

	gv, gok, _ := storage.GetGaugeMetric("g1")
	if !gok || gv != 5.5 {
		t.Errorf("expected gauge 5.5, got %f", gv)
	}
	cv, cok, _ := storage.GetCounterMetric("c1")
	if !cok || cv != 15 {
		t.Errorf("expected counter 15, got %d", cv)
	}
}

func TestUpdateMetrics(t *testing.T) {
	storage := NewMemStorageRepository()

	initial := []models.Metrics{
		{ID: "c1", MType: models.CounterMetric, Delta: utils.IntPtr(10)},
	}
	_ = storage.UpdateMetrics(initial)

	update := []models.Metrics{
		{ID: "c1", MType: models.CounterMetric, Delta: utils.IntPtr(5)},
		{ID: "g1", MType: models.GaugeMetric, Value: utils.FloatPtr(3.14)},
	}
	_ = storage.UpdateMetrics(update)

	cv, cok, _ := storage.GetCounterMetric("c1")
	if !cok || cv != 15 {
		t.Errorf("expected counter 15, got %d", cv)
	}
	gv, gok, _ := storage.GetGaugeMetric("g1")
	if !gok || gv != 3.14 {
		t.Errorf("expected gauge 3.14, got %f", gv)
	}
}
