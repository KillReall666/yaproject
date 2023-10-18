package main

import (
	"context"
	"fmt"
	db "github.com/KillReall666/yaproject/internal/db"
	"net/http"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/handlers/get"
	"github.com/KillReall666/yaproject/internal/handlers/html"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/handlers/zipdata"
	logger "github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadServerConfig()
	log, err1 := logger.InitLogger()
	if err1 != nil {
		panic("cannot initialize zap")
	}
	store := storage.NewMemStorage()
	fileWriter := fileutil.NewFileIo(cfg, store, log)

	var flag = true
	if cfg.DefaultDBConnStr == "" {
		flag = false
		log.LogInfo("Metric storage switched to memory. The database is not connected.")
	}

	db, conn, err := db.GetDB(cfg.DefaultDBConnStr)
	if err != nil {
		log.LogInfo("Database not loaded.", err)
	}

	if flag {
		err = db.CreateMetricsTable(conn)
		if err != nil {
			log.LogInfo("error creating table.", err)
		}
		defer conn.Close(context.Background())
		log.LogInfo("Database connection established.")
	}

	app := service.NewService(flag, log, fileWriter, db, store)

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

	app.LogInfo("starting http server to serve metrics on port", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		app.LogInfo("server is down:", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
