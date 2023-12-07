package main

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers/html"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/handlers/getmetrics"
	"github.com/KillReall666/yaproject/internal/handlers/zipdata"
	"github.com/KillReall666/yaproject/internal/hashmiddleware"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/servicemetric"
	"github.com/KillReall666/yaproject/internal/storage"
)

func main() {
	log, err := logger.InitLogger()
	if err != nil {
		panic("cannot initialize zap")
	}

	cfg, err := config.LoadForServer()
	if err != nil {
		log.LogInfo("cfg:", err)
	}

	store, err := storage.NewStore(cfg, log)
	if err != nil {
		log.LogInfo(err)
	}

	app := servicemetric.NewService(log, store)

	r := chi.NewRouter()
	r.Use(log.MyLogger)
	r.Use(zipdata.GzipMiddleware)
	r.Use(hashmiddleware.NewHashMiddleware(cfg.HashKey))

	r.Post("/update/*", update.NewUpdateHandler(app, log, cfg).UpdateMetrics)
	r.Post("/update/", update.NewUpdateHandler(app, log, cfg).UpdateJSONMetrics)
	r.Post("/value/", getmetrics.NewHandler(app, cfg).GetMetricsJSON)
	r.Post("/updates/", update.NewBatchHandler(app, log, cfg).BatchUpdateMetrics)

	r.Get("/value/*", getmetrics.NewHandler(app, cfg).Metrics)
	r.Get("/", html.NewHTMLHandler(app).HTMLOutput)
	r.Get("/ping", getmetrics.NewCheckDBStatusHandler(app, log).DBStatusCheck)

	log.LogInfo("starting http server to serve metrics on port", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.LogInfo("server is down:", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
