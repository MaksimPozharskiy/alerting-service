package metrics

import (
	"alerting-service/internal/models"
	"encoding/json"
	"io"
	"os"
	"sync"
)

type BackupController struct {
	file    io.ReadWriteCloser
	encoder *json.Encoder
	decoder *json.Decoder
	mutex   sync.Mutex
}

func NewBackupController(filename string) (*BackupController, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	return &BackupController{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (b *BackupController) WriteMetric(metric models.Metrics) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.encoder.Encode(&metric)
}

func (b *BackupController) WriteMetrics(metrics []models.Metrics) error {
	for _, metric := range metrics {
		err := b.WriteMetric(metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BackupController) ReadMetrics() ([]models.Metrics, error) {
	allMetrics := []models.Metrics{}

	for b.decoder.More() {
		metric, err := b.ReadMetric()
		if err != nil {
			return nil, err
		}

		allMetrics = append(allMetrics, metric)
	}

	return allMetrics, nil
}

func (b *BackupController) ReadMetric() (models.Metrics, error) {
	metric := models.Metrics{}
	if err := b.decoder.Decode(&metric); err != nil && err != io.EOF {
		return models.Metrics{}, err
	}

	return metric, nil
}
