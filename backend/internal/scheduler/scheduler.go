package scheduler

import (
	"sync"
	"time"

	"github.com/PingTower/internal/metrics"

	"github.com/PingTower/internal/entities"
)

type Scheduler struct {
	mu sync.Mutex
	// список задач
	jobs map[string]*entities.Job
	// graceful shutdown
	quit    chan struct{}
	wg      sync.WaitGroup
	metrics *metrics.SchedulerMetrics
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			s.runDueJobs()
		case <-s.quit:
			ticker.Stop()
			s.wg.Wait()
			return
		}
	}
}

func (s *Scheduler) runDueJobs() {
	now := time.Now()

	for _, job := range s.jobs {
		if now.After(job.NextRun) {
			s.wg.Add(1)

			go func(j *entities.Job) {
				defer s.wg.Done()
				start := time.Now()
				err := j.Handler(*j)
				duration := time.Since(start)

				if err != nil {
					s.metrics.RecordJobExec(duration, false)
				} else {
					s.metrics.RecordJobExec(duration, true)
				}

				j.LastRun = now
				j.NextRun = now.Add(j.Interval)
			}(job)
		}
	}
}
