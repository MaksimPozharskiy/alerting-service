package utils

import (
	"slices"
	"strconv"
	"strings"

	"alerting-service/internal/models"
	v "alerting-service/internal/validation"
)

func ParseUpdateMetricURL(url string) (models.Metrics, error) {
	var m models.Metrics

	urlData := strings.Split(url, "/")

	if len(urlData) != v.ValidCountUpdateURLParts {
		return m, v.ErrMetricNotFound
	}

	metricType := urlData[2]
	if !slices.Contains(v.ValidMetricTypes, metricType) {
		return m, v.ErrInvalidMetricType
	}

	m.MType = metricType
	m.ID = urlData[3]
	if m.MType == "gauge" {
		value, err := strconv.ParseFloat(urlData[4], 64)
		if err != nil {
			return m, v.ErrInvalidMetricValue
		}
		m.Value = &value
	} else {
		value, err := strconv.Atoi(urlData[4])
		if err != nil {
			return m, v.ErrInvalidMetricValue
		}
		val := int64(value)
		m.Delta = &val
	}

	return m, nil
}

func ParseGetMetricURL(url string) (models.Metrics, error) {
	var m models.Metrics

	urlData := strings.Split(url, "/")

	if len(urlData) != v.ValidCountGetURLParts {
		return m, v.ErrMetricNotFound
	}

	metricType := urlData[2]
	if !slices.Contains(v.ValidMetricTypes, metricType) {
		return m, v.ErrInvalidMetricType
	}

	m.MType = metricType
	m.ID = urlData[3]

	return m, nil
}
