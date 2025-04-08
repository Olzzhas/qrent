package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.Handler(http.MethodGet, "/debug", expvar.Handler())

	// Organization routes.
	router.HandlerFunc(http.MethodGet, "/v1/organizations", app.ListOrganizationHandler)
	router.HandlerFunc(http.MethodPost, "/v1/organizations", app.CreateOrganizationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/organizations/:id", app.GetOrganizationHandler)
	router.HandlerFunc(http.MethodPut, "/v1/organizations/:id", app.UpdateOrganizationHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/organizations/:id", app.DeleteOrganizationHandler)

	// Powerbank routes.
	router.HandlerFunc(http.MethodGet, "/v1/powerbanks", app.ListPowerbankHandler)
	router.HandlerFunc(http.MethodPost, "/v1/powerbanks", app.CreatePowerbankHandler)
	router.HandlerFunc(http.MethodGet, "/v1/powerbanks/:id", app.GetPowerbankHandler)
	router.HandlerFunc(http.MethodPut, "/v1/powerbanks/:id", app.UpdatePowerbankHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/powerbanks/:id", app.DeletePowerbankHandler)

	// Station routes.
	router.HandlerFunc(http.MethodGet, "/v1/stations", app.ListStationHandler)
	router.HandlerFunc(http.MethodPost, "/v1/stations", app.CreateStationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/stations/:id", app.GetStationHandler)
	router.HandlerFunc(http.MethodPut, "/v1/stations/:id", app.UpdateStationHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/stations/:id", app.DeleteStationHandler)

	return app.metrics(app.recoverPanic(app.rateLimit(router)))
}
