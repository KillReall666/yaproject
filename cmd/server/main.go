package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers/get"
	"github.com/KillReall666/yaproject/internal/handlers/html"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	store := storage.NewMemStorage()
	serv := service.NewService(store)

	getHandler := get.NewGetHandler(serv)
	updateHandler := update.NewUpdateHandler(serv)
	htmlHandler := html.NewHtmlHandler(serv)

	cfg := config.LoadServerConfig()

	r := chi.NewRouter()

	r.Post("/update/*", updateHandler.UpdateMetrics)
	r.Get("/value/*", getHandler.GetMetrics)
	r.HandleFunc("/", htmlHandler.HTMLOutput)

	log.Printf("Starting http server to serve metricss at port%s ", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Printf("server is down: %v", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
