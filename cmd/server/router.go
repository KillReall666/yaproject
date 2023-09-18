package main

import (
	"github.com/go-chi/chi/v5"
)

func MyNewRouter() (router chi.Router) {
	router = chi.NewRouter()
	router.Post("/update/*", updateHandler.UpdateMetrics)
	router.Get("/value/*", updateHandler.GetMetrics)
	router.HandleFunc("/", updateHandler.HTMLOutput)
	return router
}
