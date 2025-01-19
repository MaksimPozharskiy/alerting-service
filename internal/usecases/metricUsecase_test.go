package usecases

import (
	"alerting-service/internal/domain"
	repositories "alerting-service/internal/repository"
	v "alerting-service/internal/validation"
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

	tests := []struct {
		name    string
		wantErr error
		metric  domain.Metric
	}{
		{
			name:    "new gauge metric usecase test",
			wantErr: nil,
			metric: domain.Metric{
				Type:  "gauge",
				Name:  "temp",
				Value: "25",
			},
		},
		{
			name:    "new counter metric usecase test",
			wantErr: nil,
			metric: domain.Metric{
				Type:  "counter",
				Name:  "temp",
				Value: "25",
			},
		},
		{
			name:    "invalid metric usecase with without nums test",
			wantErr: v.ErrInvalidMetricValue,
			metric: domain.Metric{
				Type:  "counter",
				Name:  "temp",
				Value: "adsadsd",
			},
		},
		{
			name:    "invalid metric usecase with nums test",
			wantErr: v.ErrInvalidMetricValue,
			metric: domain.Metric{
				Type:  "counter",
				Name:  "temp",
				Value: "55fff",
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
