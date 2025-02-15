package agent

import (
	"alerting-service/internal/config"
	"alerting-service/internal/models"
	sign "alerting-service/internal/signature"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var pollCount int

type stats map[string]float64

func RuntimeAgent(client *http.Client) {
	conf := config.GetConfig()

	memStat := &runtime.MemStats{}

	runtime.ReadMemStats(memStat)

	var stats = make(map[string]float64)

	go func() {
		for {
			time.Sleep(time.Duration(conf.PollInterval) * time.Second)
			pollCount++
			runtime.ReadMemStats(memStat)
			getMemStatData(memStat, stats)
		}
	}()

	for {
		time.Sleep(time.Duration(conf.ReportInterval) * time.Second)
		stats["RandomValue"] = rand.Float64()
		sendMetrics(client, stats, conf)
	}
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

func sendMetrics(client *http.Client, stats stats, conf *config.Config) {
	for key, val := range stats {
		metric := models.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &val,
		}

		sendGaugeMetric(client, metric, conf)
	}

	pollCount := int64(pollCount)
	metric := models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	sendCounterMetric(client, metric, conf)
}

func sendGaugeMetric(client *http.Client, metric models.Metrics, conf *config.Config) {
	url := fmt.Sprintf("http://%s/update/gauge/%s/%f", conf.RunAddr, metric.ID, *metric.Value)

	body, err := json.Marshal(metric)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	_, err = io.Copy(g, strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	if conf.HashKey != "" {
		signature := sign.GetHash(conf.HashKey)
		req.Header.Set(sign.HashSHA256, signature)
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()
}

func sendCounterMetric(client *http.Client, metric models.Metrics, conf *config.Config) {
	url := fmt.Sprintf("http://%s/update/counter/%s/%d", conf.RunAddr, metric.ID, *metric.Delta)

	body, err := json.Marshal(metric)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	_, err = io.Copy(g, strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	if conf.HashKey != "" {
		signature := sign.GetHash(conf.HashKey)
		req.Header.Set(sign.HashSHA256, signature)
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()
}
