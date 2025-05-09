package main

import (
	"errors"
	"github.com/olzzhas/qrent/pkg/validator"
	"net/http"

	"github.com/olzzhas/qrent/internal/data"
)

// GetStationHandler godoc
// @Summary Получает станцию по ID
// @Description Возвращает станцию по идентификатору
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} StationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /stations/{id} [get]
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

// CreateStationHandler godoc
// @Summary Создаёт новую станцию
// @Description Создаёт станцию, привязанную к организации (org_id)
// @Tags stations
// @Accept json
// @Produce json
// @Param station body CreateStationRequest true "Station Data"
// @Success 201 {object} StationResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /stations [post]
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
		if errors.As(data.ErrInvalidForeignKey, &err) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"station": station}
	if err := app.writeJSON(w, http.StatusCreated, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// UpdateStationHandler godoc
// @Summary Обновляет станцию по ID
// @Description Обновляет данные станции. Обновляются только переданные поля.
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param station body UpdateStationRequest true "Station Data"
// @Success 200 {object} StationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /stations/{id} [put]
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
		if errors.As(data.ErrInvalidForeignKey, &err) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"station": station}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// DeleteStationHandler godoc
// @Summary Удаляет станцию по ID
// @Description Удаляет станцию с заданным идентификатором
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /stations/{id} [delete]
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

// ListStationHandler godoc
// @Summary Возвращает список станций
// @Description Возвращает все станции
// @Tags stations
// @Accept json
// @Produce json
// @Success 200 {object} StationListResponse
// @Failure 500 {object} map[string]string
// @Router /stations [get]
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
