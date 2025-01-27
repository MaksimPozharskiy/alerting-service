package metrics

import (
	"alerting-service/internal/models"
	"encoding/json"
	"io"
	"os"
)

type BackupController struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
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
	file, err := b.file.Stat()
	if err != nil {
		return err
	}

	isFirstMetric := file.Size() == 1

	if !isFirstMetric {
		_, err := b.file.WriteString(",")
		if err != nil {
			return err
		}
	}

	return b.encoder.Encode(&metric)
}

func (b *BackupController) WriteMetrics(metrics []models.Metrics) error {
	_, err := b.file.WriteString("[")
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		err := b.WriteMetric(metric)
		if err != nil {
			return err
		}
	}

	_, err = b.file.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (b *BackupController) ReadMetric() ([]models.Metrics, error) {
	allMetricss := []models.Metrics{}
	if err := b.decoder.Decode(&allMetricss); err != nil && err != io.EOF {
		return nil, err
	}

	return allMetricss, nil
}
