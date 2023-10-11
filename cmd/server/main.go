package main

import (
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
	fileWriterCfg := config.LoadFileIoConf()

	log, err1 := logger.InitLogger()

	if err1 != nil {
		panic("cannot initialize zap")
	}

	store := storage.NewMemStorage()
	fileWriter := fileutil.NewFileIo(fileWriterCfg, store, log)

	db, err := db.GetDB()
	if err != nil {
		log.LogInfo("db not loaded!", err)
	}

	app := service.NewService(store, log, fileWriter, db)

	getHandler := get.NewGetHandler(app)
	updateHandler := update.NewUpdateHandler(app, log)
	htmlHandler := html.NewHTMLHandler(app)
	checkConnHandler := get.NewCheckDbStatusHandler(app, log)

	fileWriter.Run()

	r := chi.NewRouter()
	r.Use(log.MyLogger)
	r.Use(zipdata.GzipMiddleware)

	r.Post("/update/*", updateHandler.UpdateMetrics)
	r.Post("/update/", updateHandler.UpdateJSONMetrics)
	r.Post("/value/", getHandler.GetMetricsJSON)

	r.Get("/value/*", getHandler.GetMetrics)
	r.Get("/", htmlHandler.HTMLOutput)
	r.Get("/ping", checkConnHandler.DbStatusCheck)

	app.LogInfo("starting http server to serve metrics on port", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		//log.Printf("server is down: %v", err)
		app.LogInfo("server is down:", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
