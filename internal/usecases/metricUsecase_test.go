package usecases

import (
	"alerting-service/internal/models"
	"alerting-service/internal/repository"
	"alerting-service/internal/utils"
	v "alerting-service/internal/validation"
	"reflect"
	"testing"
)

func TestNewMetricUsecase(t *testing.T) {
	rep := repository.NewMemStorageRepository()

	tests := []struct {
		name string
		want MetricUsecase
	}{
		{
			name: "new metric usecase test",
			want: &MetricUsecaseImpl{storageRepository: rep},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewMetricUsecase(rep)

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}

func TestMetricDataProcessing(t *testing.T) {
	usecase := NewMetricUsecase(repository.NewMemStorageRepository())
	testValue := 25.5
	testValuePtr := &testValue
	testDelta := int64(25)
	testDeltaPtr := &testDelta

	tests := []struct {
		name    string
		wantErr error
		metric  models.Metrics
	}{
		{
			name:    "new gauge metric usecase test",
			wantErr: nil,
			metric: models.Metrics{
				MType: "gauge",
				ID:    "temp",
				Value: testValuePtr,
			},
		},
		{
			name:    "new counter metric usecase test",
			wantErr: nil,
			metric: models.Metrics{
				MType: "counter",
				ID:    "temp",
				Delta: testDeltaPtr,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wantErr := usecase.MetricDataProcessing(test.metric)

			if wantErr != test.wantErr {
				t.Errorf("want: %v, got: %v", test.wantErr, wantErr)
			}
		})
	}
}

func TestGetMetricDataProcessing(t *testing.T) {
	usecase := NewMetricUsecase(repository.NewMemStorageRepository())

	// подготовка данных
	_ = usecase.MetricDataProcessing(models.Metrics{
		MType: "gauge",
		ID:    "gauge1",
		Value: utils.FloatPtr(10.1),
	})
	_ = usecase.MetricDataProcessing(models.Metrics{
		MType: "counter",
		ID:    "counter1",
		Delta: utils.IntPtr(7),
	})

	tests := []struct {
		name    string
		input   models.Metrics
		want    float64
		wantErr error
	}{
		{
			name: "get gauge value",
			input: models.Metrics{
				MType: "gauge",
				ID:    "gauge1",
			},
			want:    10.1,
			wantErr: nil,
		},
		{
			name: "get counter value",
			input: models.Metrics{
				MType: "counter",
				ID:    "counter1",
			},
			want:    7,
			wantErr: nil,
		},
		{
			name: "get unknown gauge",
			input: models.Metrics{
				MType: "gauge",
				ID:    "unknown",
			},
			want:    0,
			wantErr: v.ErrMetricNotFound,
		},
		{
			name: "invalid type",
			input: models.Metrics{
				MType: "invalid",
				ID:    "x",
			},
			want:    0,
			wantErr: v.ErrInvalidMetricValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecase.GetMetricDataProcessing(tt.input)
			if err != tt.wantErr {
				t.Errorf("wantErr %v, got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	usecase := NewMetricUsecase(repository.NewMemStorageRepository())
	_ = usecase.MetricDataProcessing(models.Metrics{
		MType: "gauge",
		ID:    "load",
		Value: utils.FloatPtr(1.0),
	})

	metrics, err := usecase.GetMetrics()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}
	if metrics[0].ID != "load" {
		t.Errorf("unexpected metric ID: %s", metrics[0].ID)
	}
}

func TestUpdateMetrics(t *testing.T) {
	usecase := NewMetricUsecase(repository.NewMemStorageRepository())

	metrics := []models.Metrics{
		{MType: "gauge", ID: "g1", Value: utils.FloatPtr(2.2)},
		{MType: "counter", ID: "c1", Delta: utils.IntPtr(5)},
	}

	err := usecase.UpdateMetrics(metrics)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	all, _ := usecase.GetMetrics()
	if len(all) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(all))
	}
}
