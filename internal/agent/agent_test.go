package agent

import (
	"runtime"
	"testing"
)

func Test_getMemStatData(t *testing.T) {
	stats := make(map[string]float64)
	memStat := &runtime.MemStats{}
	runtime.ReadMemStats(memStat)

	getMemStatData(memStat, stats)

	if stats["Alloc"] == 0 {
		t.Error("expected non-zero Alloc")
	}
	if len(stats) < 10 {
		t.Errorf("expected at least 10 metrics, got %d", len(stats))
	}
}
