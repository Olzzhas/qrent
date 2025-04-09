package main

import (
	"github.com/olzzhas/qrent/internal/data"
	"github.com/olzzhas/qrent/pkg/validator"
	"net/http"
)

// GetOrganizationHandler godoc
// @Summary Получает организацию по ID
// @Description Возвращает организацию по переданному идентификатору
// @Tags organizations
// @Accept  json
// @Produce  json
// @Param id path int true "Organization ID"
// @Success 200 {object} OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organizations/{id} [get]
func (app *application) GetOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	org, err := app.models.Organization.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"organization": org}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// CreateOrganizationHandler godoc
// @Summary Создаёт новую организацию
// @Description Создаёт организацию с заданными name и location
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body CreateOrganizationRequest true "Organization Data"
// @Success 201 {object} OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /organizations [post]
func (app *application) CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	org := &data.Organization{
		Name:     input.Name,
		Location: input.Location,
	}

	v := validator.New()
	data.ValidateOrganization(v, org)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Organization.Insert(org); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"organization": org}
	if err := app.writeJSON(w, http.StatusCreated, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// UpdateOrganizationHandler godoc
// @Summary Обновляет организацию по ID
// @Description Обновляет данные организации. Обновляются только переданные поля.
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param organization body UpdateOrganizationRequest true "Organization Data"
// @Success 200 {object} OrganizationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /organizations/{id} [put]
func (app *application) UpdateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	org, err := app.models.Organization.Get(int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Location *string `json:"location"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		org.Name = *input.Name
	}
	if input.Location != nil {
		org.Location = *input.Location
	}

	v := validator.New()
	data.ValidateOrganization(v, org)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Organization.Update(org); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"organization": org}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// DeleteOrganizationHandler godoc
// @Summary Удаляет организацию по ID
// @Description Удаляет организацию с заданным идентификатором
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /organizations/{id} [delete]
func (app *application) DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.models.Organization.Delete(int(id)); err != nil {
		app.notFoundResponse(w, r)
		return
	}

	env := envelope{"message": "organization successfully deleted"}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// ListOrganizationHandler godoc
// @Summary Возвращает список организаций
// @Description Возвращает все организации
// @Tags organizations
// @Accept json
// @Produce json
// @Success 200 {object} OrganizationListResponse
// @Failure 500 {object} map[string]string
// @Router /organizations [get]
func (app *application) ListOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	orgs, err := app.models.Organization.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{"organizations": orgs}
	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
