package scheduler

import (
	"github.com/PingTower/internal/httprequests"
	"github.com/PingTower/internal/service"
	"sync"
	"time"

	"context"
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
	Manager *service.ServiceManager
}

func NewScheduler(manager *service.ServiceManager, m *metrics.SchedulerMetrics) *Scheduler {
	return &Scheduler{
		Jobs:    make(map[string]*entities.Job),
		Quit:    make(chan struct{}),
		Metrics: m,
		Manager: manager,
	}
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

func (s *Scheduler) refreshJobs(ctx context.Context) {
	services, err := s.Manager.GetAllServices(ctx)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, srv := range services {
		if !srv.Active {
			delete(s.Jobs, srv.ID.String())
			continue
		}

		if _, exists := s.Jobs[srv.ID.String()]; !exists {
			job := &entities.Job{
				ID:       srv.ID,
				Name:     srv.ServiceName,
				Interval: srv.Interval,
				NextRun:  time.Now(),
				Handler: func(j entities.Job) (interface{}, error) {
					return httprequests.HttpCheck(j)
				},
			}
			s.Jobs[srv.ID.String()] = job
		}
	}
}

func (s *Scheduler) runDueJobs() {
	now := time.Now()

	for _, job := range s.Jobs {
		if now.After(job.NextRun) {
			s.wg.Add(1)

			go func(j *entities.Job) {
				defer s.wg.Done()
				start := time.Now()
				_, err := j.Handler(*j)
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
