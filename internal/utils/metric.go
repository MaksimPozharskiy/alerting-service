package utils

import (
	"slices"
	"strings"

	v "alerting-service/internal/validation"
)

type Metric struct {
	Type  string
	Name  string
	Value string
}

func ParseMetricURL(url string) (Metric, error) {
	var m Metric

	urlData := strings.Split(url, "/")

	if len(urlData) != v.ValidCountURLParts {
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
