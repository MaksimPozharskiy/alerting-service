package repositories

import "strconv"

type MemStorage interface {
	GetCounterMetric(string) (int, bool)
	GetGaugeMetric(string) (float64, bool)
	UpdateGaugeMetric(string, float64)
	UpdateCounterMetric(string, int)
	GetMetrics() map[string]string
}

type MemStorageImp struct {
	gauges   map[string]float64
	counters map[string]int
}

func NewStorageRepository() MemStorage {
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

func (s *MemStorageImp) GetMetrics() map[string]string {
	allMetrics := make(map[string]string)

	for key, value := range s.gauges {
		allMetrics[key] = strconv.FormatFloat(value, 'f', -1, 64)
	}

	for key, value := range s.counters {
		allMetrics[key] = strconv.Itoa(value)
	}

	return allMetrics
}
