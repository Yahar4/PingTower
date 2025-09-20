package entities

import "github.com/google/uuid"

type CreateServiceRequest struct {
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
	Interval    int    `json:"interval"`
	Active      bool   `json:"active"`
}

type CreateServiceResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateServiceRequest struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	URL         string    `json:"url"`
	Interval    int       `json:"interval"`
	Active      bool      `json:"active"`
}
