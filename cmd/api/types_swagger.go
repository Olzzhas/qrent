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
