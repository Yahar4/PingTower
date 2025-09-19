package metrics

import (
	"sync"
	"time"
)

type SchedulerMetrics struct {
	mu sync.RWMutex

	JobsExecutedTotal int64
	JobsFailedTotal   int64
	TotalResponseTime time.Duration
	ActiveWorkers     int64
	QueueLength       int64
}

func NewSchedulerMetrics() *SchedulerMetrics {
	return &SchedulerMetrics{}
}

func (m *SchedulerMetrics) RecordJobExec(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.JobsExecutedTotal++
	m.TotalResponseTime += duration

	if !success {
		m.JobsFailedTotal++
	}
}

func (m *SchedulerMetrics) GetMetrics() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgTimeMs := float64(0)
	if m.JobsExecutedTotal > 0 {
		avgTimeMs = float64(m.TotalResponseTime.Milliseconds()) / float64(m.JobsExecutedTotal)
	}

	successRate := float64(0)
	if m.JobsExecutedTotal > 0 {
		successRate = float64(m.JobsExecutedTotal-m.JobsFailedTotal) / float64(m.JobsExecutedTotal) * 100
	}

	return map[string]float64{
		"jobs_executed_total":  float64(m.JobsExecutedTotal),
		"jobs_failed_total":    float64(m.JobsFailedTotal),
		"avg_response_time_ms": avgTimeMs,
		"active_workers":       float64(m.ActiveWorkers),
		"queue_length":         float64(m.QueueLength),
		"success_rate":         successRate,
	}
}
