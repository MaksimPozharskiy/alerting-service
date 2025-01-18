package agent

import (
	"alerting-service/internal/config"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
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
		sendMetrics(client, stats, conf.RunAddr)
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

func sendMetrics(client *http.Client, stats stats, address string) {
	for key, val := range stats {
		sendGaugeMetric(client, key, val, address)
	}

	sendCounterMetric(client, "PollCount", pollCount, address)
}

func sendGaugeMetric(client *http.Client, metricName string, metricValue float64, address string) {
	url := fmt.Sprintf("http://%s/update/gauge/%s/%f/", address, metricName, metricValue)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-type", "text/plain; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()
}

func sendCounterMetric(client *http.Client, metricName string, metricValue int, address string) {
	url := fmt.Sprintf("http://%s/update/counter/%s/%d/", address, metricName, metricValue)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-type", "text/plain; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()
}
