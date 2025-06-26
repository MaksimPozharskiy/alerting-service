package utils

import (
	"alerting-service/internal/models"
	v "alerting-service/internal/validation"
	"reflect"
	"testing"
)

func TestParseUpdateMetricURL(t *testing.T) {
	testValue := 25.5
	testValuePtr := &testValue
	testDelta := int64(25)
	testDeltaPtr := &testDelta
	tests := []struct {
		name     string
		inputURL string
		want     models.Metrics
		wantErr  error
	}{
		{
			name:     "valid gauge metric url test",
			inputURL: "/update/gauge/temperature/25.5",
			want:     models.Metrics{MType: "gauge", ID: "temperature", Value: testValuePtr},
			wantErr:  nil,
		},
		{
			name:     "valid counter metric url test",
			inputURL: "/update/counter/temperature/25",
			want:     models.Metrics{MType: "counter", ID: "temperature", Delta: testDeltaPtr},
			wantErr:  nil,
		},
		{
			name:     "invalid metric url test",
			inputURL: "/update/cter/temperature/25",
			want:     models.Metrics{},
			wantErr:  v.ErrInvalidMetricType,
		},
		{
			name:     "not found metric url test",
			inputURL: "/update/gauge/temperature",
			want:     models.Metrics{},
			wantErr:  v.ErrMetricNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want, wantErr := ParseUpdateMetricURL(test.inputURL)

			if wantErr != test.wantErr {
				t.Errorf("want error: %v, got: %v", test.wantErr, wantErr)
			}

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %+v, got: %+v", test.want, want)
			}
		})
	}
}

func TestParseGetMetricURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		want     models.Metrics
		wantErr  error
	}{
		{
			name:     "valid gauge metric get url",
			inputURL: "/value/gauge/temperature",
			want:     models.Metrics{MType: "gauge", ID: "temperature"},
			wantErr:  nil,
		},
		{
			name:     "valid counter metric get url",
			inputURL: "/value/counter/hits",
			want:     models.Metrics{MType: "counter", ID: "hits"},
			wantErr:  nil,
		},
		{
			name:     "invalid metric type",
			inputURL: "/value/ctr/hits",
			want:     models.Metrics{},
			wantErr:  v.ErrInvalidMetricType,
		},
		{
			name:     "not enough parts in url",
			inputURL: "/value/gauge",
			want:     models.Metrics{},
			wantErr:  v.ErrMetricNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseGetMetricURL(test.inputURL)

			if err != test.wantErr {
				t.Errorf("expected err: %v, got: %v", test.wantErr, err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("expected: %+v, got: %+v", test.want, got)
			}
		})
	}
}
