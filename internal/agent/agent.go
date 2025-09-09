package agent

import (
	"alerting-service/internal/config"
	"alerting-service/internal/crypto"
	"alerting-service/internal/logger"
	"alerting-service/internal/models"
	sign "alerting-service/internal/signature"
	"context"
	"crypto/rsa"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
)

var pollCount int

type stats map[string]float64

type sentMetricWorker struct {
	client    *http.Client
	conf      *config.Config
	publicKey *rsa.PublicKey
}

func RuntimeAgent(ctx context.Context, client *http.Client) {
	conf := config.GetConfig()
	memStat := &runtime.MemStats{}
	runtime.ReadMemStats(memStat)
	stats := make(map[string]float64)
	metricsChan := make(chan models.Metrics, 30)
	resultsChan := make(chan error, conf.RateLimit)

	var publicKey *rsa.PublicKey
	var err error
	if conf.CryptoKey != "" {
		publicKey, err = crypto.LoadPublicKey(conf.CryptoKey)
		if err != nil {
			logger.Log.Error("Failed to load public key, proceeding without encryption", zap.Error(err))
		}
	}

	for w := 1; w <= conf.RateLimit; w++ {
		worker := sentMetricWorker{client: client, conf: conf, publicKey: publicKey}
		go worker.sendMetric(ctx, metricsChan, resultsChan)
	}

	go func() {
		ticker := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pollCount++
				runtime.ReadMemStats(memStat)
				getMemStatData(memStat, stats)
			}
		}
	}()

	<-ctx.Done()
	close(metricsChan)
}

func getMemStatData(memStat *runtime.MemStats, stats stats) {
	stats["Alloc"] = float64(memStat.Alloc)
	stats["BuckHashSys"] = float64(memStat.BuckHashSys)
	stats["Frees"] = float64(memStat.Frees)
	stats["GCCPUFraction"] = float64(memStat.GCCPUFraction)
	stats["GCSys"] = float64(memStat.GCSys)
	stats["HeapAlloc"] = float64(memStat.HeapAlloc)
	stats["HeapIdle"] = float64(memStat.HeapIdle)
	stats["HeapInuse"] = float64(memStat.HeapInuse)
	stats["HeapObjects"] = float64(memStat.HeapObjects)
	stats["HeapReleased"] = float64(memStat.HeapReleased)
	stats["HeapSys"] = float64(memStat.HeapSys)
	stats["LastGC"] = float64(memStat.LastGC)
	stats["Lookups"] = float64(memStat.Lookups)
	stats["MCacheInuse"] = float64(memStat.MCacheInuse)
	stats["MCacheSys"] = float64(memStat.MCacheSys)
	stats["MSpanInuse"] = float64(memStat.MSpanInuse)
	stats["MSpanSys"] = float64(memStat.MSpanSys)
	stats["Mallocs"] = float64(memStat.Mallocs)
	stats["NextGC"] = float64(memStat.NextGC)
	stats["NumForcedGC"] = float64(memStat.NumForcedGC)
	stats["NumGC"] = float64(memStat.NumGC)
	stats["OtherSys"] = float64(memStat.OtherSys)
	stats["PauseTotalNs"] = float64(memStat.PauseTotalNs)
	stats["StackInuse"] = float64(memStat.StackInuse)
	stats["StackSys"] = float64(memStat.StackSys)
	stats["Sys"] = float64(memStat.Sys)
	stats["TotalAlloc"] = float64(memStat.TotalAlloc)
}

func sendMetrics(stats stats, metricChan chan models.Metrics) {
	for key, value := range stats {
		metric := models.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}

		metricChan <- metric
	}

	pollCount := int64(pollCount)
	metric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}

	metricChan <- metric
}
func (w *sentMetricWorker) sendGaugeMetric(metric models.Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Error("marshal error", zap.Error(err))
		return err
	}

	finalBody := body
	if w.publicKey != nil {
		finalBody, err = crypto.EncryptData(body, w.publicKey)
		if err != nil {
			logger.Log.Error("encryption error", zap.Error(err))
			return err
		}
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	_, err = g.Write(finalBody)
	if err != nil {
		logger.Log.Error("gzip write error", zap.Error(err))
		g.Close()
		return err
	}
	g.Close()

	url := fmt.Sprintf("http://%s/update", w.conf.RunAddr)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		logger.Log.Error("request error", zap.Error(err))
		return err
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	if w.publicKey != nil {
		req.Header.Set("Content-Encryption", "RSA")
	}

	if w.conf.HashKey != "" {
		signature := sign.GetHash(body, []byte(w.conf.HashKey))
		req.Header.Set(sign.HashSHA256, signature)
	}

	response, err := w.client.Do(req)
	if err != nil {
		logger.Log.Error("http send error", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	return nil
}

func (w *sentMetricWorker) sendCounterMetric(metric models.Metrics) error {
	body, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Error("marshal error", zap.Error(err))
		return err
	}

	finalBody := body
	if w.publicKey != nil {
		finalBody, err = crypto.EncryptData(body, w.publicKey)
		if err != nil {
			logger.Log.Error("encryption error", zap.Error(err))
			return err
		}
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	_, err = g.Write(finalBody)
	if err != nil {
		logger.Log.Error("gzip write error", zap.Error(err))
		g.Close()
		return err
	}
	g.Close()

	url := fmt.Sprintf("http://%s/update", w.conf.RunAddr)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		logger.Log.Error("request error", zap.Error(err))
		return err
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	if w.publicKey != nil {
		req.Header.Set("Content-Encryption", "RSA")
	}

	if w.conf.HashKey != "" {
		signature := sign.GetHash(body, []byte(w.conf.HashKey))
		req.Header.Set(sign.HashSHA256, signature)
	}

	response, err := w.client.Do(req)
	if err != nil {
		logger.Log.Error("http send error", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	return nil
}

func (w *sentMetricWorker) sendMetric(ctx context.Context, metricsChan <-chan models.Metrics, resultsChan chan<- error) {
	for {
		select {
		case metric, ok := <-metricsChan:
			if !ok {
				return
			}
			var err error
			switch metric.MType {
			case models.GaugeMetric:
				err = w.sendGaugeMetric(metric)
			case models.CounterMetric:
				err = w.sendCounterMetric(metric)
			}
			resultsChan <- err
		case <-ctx.Done():
			return
		}
	}
}
