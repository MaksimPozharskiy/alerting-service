package models

// CounterMetric is the type identifier for counter metrics.
const CounterMetric = "counter"

// GaugeMetric is the type identifier for gauge metrics.
const GaugeMetric = "gauge"

// Metrics defines a data structure representing a single metric.
type Metrics struct {
	ID    string   `json:"id"`              // Unique metric identifier
	MType string   `json:"type"`            // Metric type: "gauge" or "counter"
	Delta *int64   `json:"delta,omitempty"` // Metric value for counter type
	Value *float64 `json:"value,omitempty"` // Metric value for gauge type
}
