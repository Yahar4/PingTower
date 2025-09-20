package entities

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID       uuid.UUID
	Name     string
	URL      string
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
	Handler  func(Job) (interface{}, error)
}

type Service struct {
	ID          uuid.UUID
	ServiceName string
	URL         string
	Interval    time.Duration
	Active      bool
}
