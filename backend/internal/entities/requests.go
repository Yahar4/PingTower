package entities

import "github.com/google/uuid"

type CreateServiceRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	Active   bool   `json:"active"`
}

type CreateServiceResponse struct {
	ID uuid.UUID `json:"id"`
}
