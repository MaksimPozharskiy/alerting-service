package repositories

type MemStorage interface {
	GetCounterMetric(string) (int, bool)
	GetGaugesMetric(string) (float64, bool)
	UpdateGaugeMetic(string, float64)
	UpdateCounterMetic(string, int)
}

type MemStorageImp struct {
	gauges   map[string]float64
	counters map[string]int
}

func NewStorageRepository() MemStorage {
	return &MemStorageImp{gauges: map[string]float64{}, counters: map[string]int{}}
}

func (s *MemStorageImp) GetCounterMetric(key string) (int, bool) {
	if counter, ok := s.counters[key]; ok {
		return 0, ok
	} else {
		return counter, false
	}
}

func (s *MemStorageImp) GetGaugesMetric(key string) (float64, bool) {
	if gauge, ok := s.gauges[key]; ok {
		return 0, ok
	} else {
		return gauge, false
	}
}

func (s *MemStorageImp) UpdateGaugeMetic(metric string, value float64) {
	s.gauges[metric] = value
}

func (s *MemStorageImp) UpdateCounterMetic(metric string, value int) {
	s.counters[metric] += value
}
