package utils

import (
	"alerting-service/internal/domain"
	v "alerting-service/internal/validation"
	"testing"
)

func TestParseUpdateMetricURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		want     domain.Metric
		wantErr  error
	}{
		{
			name:     "valid gauge metric url test",
			inputURL: "/update/gauge/temperature/25.5",
			want:     domain.Metric{Type: "gauge", Name: "temperature", Value: "25.5"},
			wantErr:  nil,
		},
		{
			name:     "valid counter metric url test",
			inputURL: "/update/counter/temperature/25",
			want:     domain.Metric{Type: "counter", Name: "temperature", Value: "25"},
			wantErr:  nil,
		},
		{
			name:     "invalid metric url test",
			inputURL: "/update/cter/temperature/25",
			want:     domain.Metric{},
			wantErr:  v.ErrInvalidMetricType,
		},
		{
			name:     "not found metric url test",
			inputURL: "/update/gauge/temperature",
			want:     domain.Metric{},
			wantErr:  v.ErrMetricNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want, wantErr := ParseUpdateMetricURL(test.inputURL)

			if wantErr != test.wantErr {
				t.Errorf("want error: %v, got: %v", test.wantErr, wantErr)
			}

			if want.Type != test.want.Type || want.Name != test.want.Name || want.Value != test.want.Value {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}
