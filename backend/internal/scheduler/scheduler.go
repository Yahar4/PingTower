package scheduler

import (
	"sync"
	"time"

	"github.com/PingTower/internal/metrics"

	"github.com/PingTower/internal/entities"
)

// TOOD: добавить функцию получения job'ов из БД

type Scheduler struct {
	mu sync.Mutex
	// список задач
	Jobs map[string]*entities.Job
	// graceful shutdown
	Quit    chan struct{}
	wg      sync.WaitGroup
	Metrics *metrics.SchedulerMetrics
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			s.runDueJobs()
		case <-s.Quit:
			ticker.Stop()
			s.wg.Wait()
			return
		}
	}
}

func (s *Scheduler) runDueJobs() {
	// функция получения всех джобов сюда

	now := time.Now()

	for _, job := range s.Jobs {
		if now.After(job.NextRun) {
			s.wg.Add(1)

			go func(j *entities.Job) {
				defer s.wg.Done()
				start := time.Now()
				err := j.Handler(*j)
				duration := time.Since(start)

				if err != nil {
					s.Metrics.RecordJobExec(duration, false)
				} else {
					s.Metrics.RecordJobExec(duration, true)
				}

				// запрос в БД на изменения

				j.LastRun = now
				j.NextRun = now.Add(j.Interval)
			}(job)
		}
	}
}
