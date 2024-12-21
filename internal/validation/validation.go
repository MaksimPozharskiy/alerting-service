package validation

import (
	"alerting-service/internal/models"
)

var ValidMetricTypes = []string{models.CounterMetric, models.GaugeMetric}
var ValidCountURLParts = 5
