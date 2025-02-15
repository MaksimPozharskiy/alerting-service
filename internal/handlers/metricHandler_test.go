package handlers

import (
	repository "alerting-service/internal/repository"
	"alerting-service/internal/usecases"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricHandler(t *testing.T) {
	metricUsecase := usecases.NewMetricUsecase(repository.NewMemStorageRepository())

	tests := []struct {
		name string
		want *metricHandler
	}{
		{
			name: "new metric handler test",
			want: &metricHandler{metricUsecase: metricUsecase},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewMetricHandler(metricUsecase)

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}

func TestUpdateMetric(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))

	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
	}{
		{
			name:   "update gauge metric success test",
			method: http.MethodPost,
			body: `{
				"id": "1",
				"type": "gauge",
				"value": 2
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name:   "update counter metric success test",
			method: http.MethodPost,
			body: `{
	  			"id": "1",
				"type": "counter",
				"delta": 2
			}`,
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			handler.UpdateMetric(w, req)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode, "Response code didn't match expected")
		})
	}
}
