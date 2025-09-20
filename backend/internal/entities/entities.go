package entities

import (
	"time"

	"github.com/google/uuid"
)

// TODO: необходимо подумать что еще мы будем привязывать к джобу
type Job struct {
	ID       uuid.UUID
	URL      string
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
	Handler  func(job Job) error
}

type Service struct {
	ID          uuid.UUID
	ServiceName string
	URL         string
	Interval    time.Duration
	Active      bool
}
