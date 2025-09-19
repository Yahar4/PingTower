package main

import (
	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/metrics"
	"github.com/PingTower/internal/scheduler"
)

func main() {
	m := metrics.NewSchedulerMetrics()
	sch := &scheduler.Scheduler{
		// пока что map, не подключил БД
		// см. комменты в scheduler.go
		Jobs:    make(map[string]*entities.Job),
		Quit:    make(chan struct{}),
		Metrics: m,
	}
}
