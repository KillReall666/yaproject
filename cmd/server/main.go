package main

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/metrics"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	store := storage.NewMemStorage()
	metricsStore := metrics.NewGaugeMetricsStorage()
	serv := service.NewService(store, metricsStore)
	updateHandler := update.NewHandler(serv)

	router.Post("/update/*", updateHandler.PostHandle)
	router.Get("/value/*", updateHandler.GetHandle)
	router.HandleFunc("/", updateHandler.HTMLHandle)

	cfg := config.LoadServerConfig()
	fmt.Println(cfg)
	log.Printf("Starting http server to serve metricss at port%s ", cfg.Address)
	err := http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Printf("server is down: %v", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
