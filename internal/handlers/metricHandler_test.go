package handlers

import (
	"alerting-service/internal/models"
	repository "alerting-service/internal/repository"
	"alerting-service/internal/usecases"
	"encoding/json"
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

func TestGetMetric(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))

	handler.metricUsecase.MetricDataProcessing(models.Metrics{
		MType: "gauge",
		ID:    "cpu",
		Value: floatPtr(55.5),
	})

	body := `{"id": "cpu", "type": "gauge"}`
	req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.GetMetric(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var got models.Metrics
	err := json.NewDecoder(res.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, "cpu", got.ID)
	assert.Equal(t, "gauge", got.MType)
	assert.Equal(t, 55.5, *got.Value)
}

func TestGetAllMetrics(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))
	handler.metricUsecase.MetricDataProcessing(models.Metrics{MType: "gauge", ID: "cpu", Value: floatPtr(70.5)})
	handler.metricUsecase.MetricDataProcessing(models.Metrics{MType: "counter", ID: "hits", Delta: int64Ptr(5)})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.GetAllMetrics(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateMetrics(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))

	body := `[{"id":"cpu","type":"gauge","value":99.9},{"id":"req","type":"counter","delta":10}]`
	req := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdateMetrics(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetURLMetric(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))
	handler.metricUsecase.MetricDataProcessing(models.Metrics{MType: "gauge", ID: "load", Value: floatPtr(1.23)})

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/load", nil)
	w := httptest.NewRecorder()

	handler.GetURLMetric(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateURLMetric(t *testing.T) {
	handler := NewMetricHandler(usecases.NewMetricUsecase(repository.NewMemStorageRepository()))

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/temp/22.3", nil)
	w := httptest.NewRecorder()

	handler.UpdateURLMetric(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func floatPtr(v float64) *float64 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
