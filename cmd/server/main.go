package main

import (
	"fmt"

	"log"
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers/get"
	"github.com/KillReall666/yaproject/internal/handlers/html"
	"github.com/KillReall666/yaproject/internal/handlers/update"

	"github.com/KillReall666/yaproject/internal/handlers/zip_data"
	logger2 "github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	myLog, err1 := logger2.InitLogger()
	if err1 != nil {
		panic("cannot initialize zap")
	}

	store := storage.NewMemStorage()
	serv := service.NewService(store)

	getHandler := get.NewGetHandler(serv)
	updateHandler := update.NewUpdateHandler(serv)
	htmlHandler := html.NewHTMLHandler(serv)

	cfg := config.LoadServerConfig()

	r := chi.NewRouter()

	r.Use(myLog.MyLogger)
	r.Use(zip_data.GzipMiddleware)

	r.Post("/update/*", updateHandler.UpdateMetrics)
	r.Post("/update/", updateHandler.UpdateJSONMetrics)
	r.Post("/value/", getHandler.GetMetricsJSON)

	r.Get("/value/*", getHandler.GetMetrics)
	r.Get("/", htmlHandler.HTMLOutput)


	log.Printf("Starting http server to serve metricss at port%s ", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Printf("server is down: %v", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
