package repositories

import (
	"alerting-service/internal/models"
)

type MemStorageImp struct {
	gauges   map[string]float64
	counters map[string]int
}

func NewStorageRepository() *MemStorageImp {
	return &MemStorageImp{gauges: map[string]float64{}, counters: map[string]int{}}
}

func (s *MemStorageImp) GetCounterMetric(key string) (int, bool) {
	if counter, ok := s.counters[key]; !ok {
		return 0, ok
	} else {
		return counter, ok
	}
}

func (s *MemStorageImp) GetGaugeMetric(key string) (float64, bool) {
	if gauge, ok := s.gauges[key]; !ok {
		return 0.0, ok
	} else {
		return gauge, ok
	}
}

func (s *MemStorageImp) UpdateGaugeMetric(metricName string, value float64) {
	s.gauges[metricName] = value
}

func (s *MemStorageImp) UpdateCounterMetric(metricName string, value int) {
	s.counters[metricName] += value
}

func (s *MemStorageImp) GetMetrics() []models.Metrics {
	allMetrics := []models.Metrics{}

	for key, value := range s.gauges {
		metric := models.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}

		allMetrics = append(allMetrics, metric)
	}

	for key, value := range s.counters {
		metric := models.Metrics{
			ID:    key,
			MType: "counter",
		}

		val := int64(value)
		metric.Delta = &val

		allMetrics = append(allMetrics, metric)
	}

	return allMetrics
}

func (s *MemStorageImp) SetMetrics(allMetrics []models.Metrics) {
	for _, metric := range allMetrics {
		if metric.MType == models.GaugeMetric {
			s.gauges[metric.ID] = *metric.Value
		} else {
			s.counters[metric.ID] = int(*metric.Delta)
		}
	}
}
