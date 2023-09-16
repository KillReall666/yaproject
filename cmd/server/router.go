package main

import (
	"github.com/go-chi/chi/v5"
)

func MyNewRouter() (router chi.Router) {
	router = chi.NewRouter()
	router.Post("/update/*", UpdateHandler.UpdateMetrics)
	router.Get("/value/*", UpdateHandler.GetMetrics)
	router.HandleFunc("/", UpdateHandler.HTMLHandle)
	return router
}
