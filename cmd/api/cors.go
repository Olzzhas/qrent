package main

import (
	"github.com/rs/cors"
	"net/http"
)

func CorsSettings() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodPut,
		},
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		OptionsPassthrough: true,
		ExposedHeaders: []string{
			"Content-Type",
		},
		Debug: false,
	})
	return c
}
