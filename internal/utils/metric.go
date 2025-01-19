package utils

import (
	"slices"
	"strings"

	"alerting-service/internal/domain"
	v "alerting-service/internal/validation"
)

func ParseUpdateMetricURL(url string) (domain.Metric, error) {
	var m domain.Metric

	urlData := strings.Split(url, "/")

	if len(urlData) != v.ValidCountUpdateURLParts {
		return m, v.ErrMetricNotFound
	}

	metricType := urlData[2]
	if !slices.Contains(v.ValidMetricTypes, metricType) {
		return m, v.ErrInvalidMetricType
	}

	m.Type = metricType
	m.Name = urlData[3]
	m.Value = urlData[4]

	return m, nil
}

func ParseGetMetricURL(url string) (domain.Metric, error) {
	var m domain.Metric

	urlData := strings.Split(url, "/")

	if len(urlData) != v.ValidCountGetURLParts {
		return m, v.ErrMetricNotFound
	}

	metricType := urlData[2]
	if !slices.Contains(v.ValidMetricTypes, metricType) {
		return m, v.ErrInvalidMetricType
	}

	m.Type = metricType
	m.Name = urlData[3]

	return m, nil
}
