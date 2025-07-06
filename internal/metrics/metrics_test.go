package metrics

import (
	"alerting-service/internal/models"
	"alerting-service/internal/utils"
	"encoding/json"
	"os"
	"testing"
)

func TestWriteAndReadSingleMetric(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_backup_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	bc, err := NewBackupController(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to create BackupController: %v", err)
	}
	defer bc.file.Close()

	metric := models.Metrics{
		ID:    "test_gauge",
		MType: "gauge",
		Value: utils.FloatPtr(3.14),
	}

	if err = bc.WriteMetric(metric); err != nil {
		t.Fatalf("failed to write metric: %v", err)
	}

	if f, ok := bc.file.(*os.File); ok {
		f.Seek(0, 0)
		bc.decoder = json.NewDecoder(f)
	} else {
		t.Fatal("bc.file is not *os.File, cannot seek")
	}

	read, err := bc.ReadMetrics()
	if err != nil {
		t.Fatalf("failed to read metrics: %v", err)
	}

	if len(read) != 1 || read[0].ID != "test_gauge" || *read[0].Value != 3.14 {
		t.Errorf("unexpected result: %+v", read)
	}
}

func TestWriteAndReadMultipleMetrics(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "metrics_backup_multi_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	bc, err := NewBackupController(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to create BackupController: %v", err)
	}
	defer bc.file.Close()

	metricsToWrite := []models.Metrics{
		{ID: "g1", MType: "gauge", Value: utils.FloatPtr(1.23)},
		{ID: "c1", MType: "counter", Delta: utils.IntPtr(10)},
	}

	if err = bc.WriteMetrics(metricsToWrite); err != nil {
		t.Fatalf("failed to write metrics: %v", err)
	}

	if f, ok := bc.file.(*os.File); ok {
		f.Seek(0, 0)
		bc.decoder = json.NewDecoder(f)
	} else {
		t.Fatal("bc.file is not *os.File, cannot seek")
	}

	read, err := bc.ReadMetrics()
	if err != nil {
		t.Fatalf("failed to read metrics: %v", err)
	}

	if len(read) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(read))
	}

	if read[0].ID != "g1" || *read[0].Value != 1.23 {
		t.Errorf("unexpected first metric: %+v", read[0])
	}
	if read[1].ID != "c1" || *read[1].Delta != 10 {
		t.Errorf("unexpected second metric: %+v", read[1])
	}
}
