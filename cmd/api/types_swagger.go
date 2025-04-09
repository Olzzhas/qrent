package main

import "github.com/olzzhas/qrent/internal/data"

// OrganizationResponse описывает структуру ответа с одним объектом Organization.
type OrganizationResponse struct {
	Organization data.Organization `json:"organization"`
}

// OrganizationListResponse описывает структуру ответа со списком объектов Organization.
type OrganizationListResponse struct {
	Organizations []data.Organization `json:"organizations"`
}

// MessageResponse описывает структуру ответа с полем message (например, при удалении).
type MessageResponse struct {
	Message string `json:"message"`
}

// PowerbankResponse описывает ответ с одним Powerbank
// swagger:model
type PowerbankResponse struct {
	Powerbank data.Powerbank `json:"powerbank"`
}

// PowerbankListResponse описывает ответ со списком Powerbank
// swagger:model
type PowerbankListResponse struct {
	Powerbanks []data.Powerbank `json:"powerbanks"`
}

// StationResponse описывает ответ с одной Station
// swagger:model
type StationResponse struct {
	Station data.Station `json:"station"`
}

// StationListResponse описывает ответ со списком Station
// swagger:model
type StationListResponse struct {
	Stations []data.Station `json:"stations"`
}

// requests

// Organization

// CreateOrganizationRequest описывает тело запроса для создания организации.
type CreateOrganizationRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// UpdateOrganizationRequest описывает тело запроса для обновления организации.
type UpdateOrganizationRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
}

// Powerbank

// CreatePowerbankRequest описывает тело запроса для создания повербанка.
type CreatePowerbankRequest struct {
	CurrentStationID int    `json:"current_station_id"`
	Status           string `json:"status"`
}

// UpdatePowerbankRequest описывает тело запроса для обновления повербанка.
type UpdatePowerbankRequest struct {
	CurrentStationID *int    `json:"current_station_id"`
	Status           *string `json:"status"`
}

// Station

// CreateStationRequest описывает тело запроса для создания станции.
type CreateStationRequest struct {
	OrgID int `json:"org_id"`
}

// UpdateStationRequest описывает тело запроса для обновления станции.
type UpdateStationRequest struct {
	OrgID *int `json:"org_id"`
}

// ErrorResponse описывает ответ с ошибкой.
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}
