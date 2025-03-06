package metrics

import (
	"testing"
)

func TestRecordMetric(t *testing.T) {
	m := NewMetrics()
	m.RecordMetric("test_proxy", 200, 100)

	if m.GetSuccessCount("test_proxy") != 1 {
		t.Errorf("Expected success count to be 1, got %d", m.GetSuccessCount("test_proxy"))
	}
	if m.GetAverageLatency("test_proxy") != 100 {
		t.Errorf("Expected average latency to be 100, got %d", m.GetAverageLatency("test_proxy"))
	}
}

func TestRecordFailure(t *testing.T) {
	m := NewMetrics()
	m.RecordMetric("test_proxy", 500, 200)

	if m.GetFailureCount("test_proxy") != 1 {
		t.Errorf("Expected failure count to be 1, got %d", m.GetFailureCount("test_proxy"))
	}
	if m.GetAverageLatency("test_proxy") != 200 {
		t.Errorf("Expected average latency to be 200, got %d", m.GetAverageLatency("test_proxy"))
	}
}

func TestGetNonExistentProxyMetrics(t *testing.T) {
	m := NewMetrics()

	if m.GetSuccessCount("non_existent_proxy") != 0 {
		t.Errorf("Expected success count to be 0 for non-existent proxy")
	}
	if m.GetFailureCount("non_existent_proxy") != 0 {
		t.Errorf("Expected failure count to be 0 for non-existent proxy")
	}
}
