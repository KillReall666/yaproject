package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/handlers/get"
	"github.com/KillReall666/yaproject/internal/handlers/html"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/handlers/zipdata"
	logger "github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/KillReall666/yaproject/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
)

func main() {
	log, err1 := logger.InitLogger()
	if err1 != nil {
		panic("cannot initialize zap")
	}

	cfg, useDB, err := config.LoadServerConfig()
	if err != nil {
		log.LogInfo("config not loaded: ", err)
	}
	store := storage.NewMemStorage()
	fileWriter := fileutil.NewFileIo(cfg, store, log)

	db, conn, err := postgres.NewDB(cfg.DefaultDBConnStr)
	if err != nil {
		log.LogInfo("Database not loaded: ", err)
	}

	app := NewService(useDB, log, fileWriter, db, store)

	if useDB {
		err = db.CreateMetricsTable(conn)
		if err != nil {
			log.LogInfo("error creating table.", err)
		}
		defer conn.Close(context.Background())
		log.LogInfo("Database connection established.")
	}

	getHandler := get.NewGetHandler(app, cfg)
	updateHandler := update.NewUpdateHandler(app, log, cfg)
	htmlHandler := html.NewHTMLHandler(app)
	checkConnHandler := get.NewCheckDBStatusHandler(app, log)
	packHandler := update.NewPackHandler(app, log, cfg)
	fileWriter.Run()

	r := chi.NewRouter()
	r.Use(log.MyLogger)
	r.Use(zipdata.GzipMiddleware)

	r.Post("/update/*", updateHandler.UpdateMetrics)
	r.Post("/update/", updateHandler.UpdateJSONMetrics)
	r.Post("/value/", getHandler.GetMetricsJSON)
	r.Post("/updates/", packHandler.PackUpdateMetrics)

	r.Get("/value/*", getHandler.GetMetrics)
	r.Get("/", htmlHandler.HTMLOutput)
	r.Get("/ping", checkConnHandler.DBStatusCheck)

	log.LogInfo("starting http server to serve metrics on port", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.LogInfo("server is down:", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
