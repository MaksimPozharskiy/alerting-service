package handlers

import (
	repositories "alerting-service/internal/repository"
	"alerting-service/internal/usecases"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricHandler(t *testing.T) {
	metricUsecase := usecases.NewMetricUsecase(repositories.NewStorageRepository())

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
	handler := NewMetricHandler(usecases.NewMetricUsecase(repositories.NewStorageRepository()))

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "postive update gauge test",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/temp/25",
		},
		{
			name: "postive update counter test",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/temp/25",
		},
		{
			name: "incorrect metric typ test",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/test/temp/25",
		},
		{
			name: "incorrect metric value test",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/test/temp/",
		},
		{
			name: "incorrect url path value test",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/test/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)

			w := httptest.NewRecorder()

			handler.UpdateMetric(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
