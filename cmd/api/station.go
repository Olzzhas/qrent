package main

import (
	"github.com/olzzhas/qrent/pkg/validator"
	"net/http"

	"github.com/olzzhas/qrent/internal/data"
)

func (app *application) GetStationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	station, err := app.models.Station.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"station": station}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CreateStationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrgID int `json:"org_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	station := &data.Station{
		OrgID: input.OrgID,
	}

	v := validator.New()
	data.ValidateStation(v, station)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Station.Insert(station); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"station": station}
	if err := app.writeJSON(w, http.StatusCreated, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdateStationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	station, err := app.models.Station.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		OrgID *int `json:"org_id"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.OrgID != nil {
		station.OrgID = *input.OrgID
	}

	v := validator.New()
	data.ValidateStation(v, station)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Station.Update(station); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"station": station}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteStationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.models.Station.Delete(int(id)); err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"message": "station successfully deleted"}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ListStationHandler(w http.ResponseWriter, r *http.Request) {
	stations, err := app.models.Station.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"stations": stations}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
