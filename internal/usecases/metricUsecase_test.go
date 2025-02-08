package usecases

import (
	"alerting-service/internal/models"
	repositories "alerting-service/internal/repository"
	"reflect"
	"testing"
)

func TestNewMetricUsecase(t *testing.T) {
	rep := repositories.NewStorageRepository()

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
	usecase := NewMetricUsecase(repositories.NewStorageRepository())
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
				MType:  "counter",
				ID:  "temp",
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
