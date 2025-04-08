package main

import (
	"errors"
	"github.com/olzzhas/qrent/internal/data"
	"github.com/olzzhas/qrent/pkg/validator"
	"net/http"
)

func (app *application) GetPowerbankHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p, err := app.models.Powerbank.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"powerbank": p}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CreatePowerbankHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CurrentStationID int    `json:"current_station_id"`
		Status           string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p := &data.Powerbank{
		CurrentStationID: input.CurrentStationID,
		Status:           data.PowerbankStatus(input.Status),
	}

	if !p.Status.IsValid() {
		app.badRequestResponse(w, r, errors.New("invalid powerbank status"))
		return
	}

	v := validator.New()
	data.ValidatePowerbank(v, p)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Вставка нового повербанка в базу.
	if err := app.models.Powerbank.Insert(p); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"powerbank": p}
	if err := app.writeJSON(w, http.StatusCreated, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdatePowerbankHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p, err := app.models.Powerbank.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		CurrentStationID *int    `json:"current_station_id"`
		Status           *string `json:"status"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.CurrentStationID != nil {
		p.CurrentStationID = *input.CurrentStationID
	}

	if input.Status != nil {
		newStatus := data.PowerbankStatus(*input.Status)
		if !newStatus.IsValid() {
			app.badRequestResponse(w, r, errors.New("invalid powerbank status"))
			return
		}
		p.Status = newStatus
	}

	v := validator.New()
	data.ValidatePowerbank(v, p)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Powerbank.Update(p); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"powerbank": p}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeletePowerbankHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.models.Powerbank.Delete(int(id)); err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"message": "powerbank successfully deleted"}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ListPowerbankHandler(w http.ResponseWriter, r *http.Request) {
	powerbanks, err := app.models.Powerbank.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"powerbanks": powerbanks}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
