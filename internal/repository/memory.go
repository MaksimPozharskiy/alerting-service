package repository

import (
	"alerting-service/internal/models"
	"sync"
)

type MemStorageImp struct {
	gauges   map[string]float64
	counters map[string]int
	mu       sync.Mutex
}

func NewMemStorageRepository() StorageRepository {
	return &MemStorageImp{gauges: map[string]float64{}, counters: map[string]int{}}
}

func (s *MemStorageImp) GetCounterMetric(key string) (int, bool, error) {
	if counter, ok := s.counters[key]; !ok {
		return 0, ok, nil
	} else {
		return counter, ok, nil
	}
}

func (s *MemStorageImp) GetGaugeMetric(key string) (float64, bool, error) {
	if gauge, ok := s.gauges[key]; !ok {
		return 0.0, ok, nil
	} else {
		return gauge, ok, nil
	}
}

func (s *MemStorageImp) UpdateGaugeMetric(metricName string, value float64) error {
	s.mu.Lock()
	s.gauges[metricName] = value
	s.mu.Unlock()
	return nil
}

func (s *MemStorageImp) UpdateCounterMetric(metricName string, value int) error {
	s.mu.Lock()
	s.counters[metricName] += value
	s.mu.Unlock()
	return nil
}

func (s *MemStorageImp) GetMetrics() ([]models.Metrics, error) {
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

	return allMetrics, nil
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

func (s *MemStorageImp) UpdateMetrics(metrics []models.Metrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, metric := range metrics {
		if metric.MType == models.GaugeMetric && metric.Value != nil {
			s.gauges[metric.ID] = *metric.Value
		}
		if metric.MType == models.CounterMetric && metric.Delta != nil {
			s.counters[metric.ID] += int(*metric.Delta)
		}
	}
	return nil
}
