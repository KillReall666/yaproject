package main

import (
	"fmt"

	logger2 "github.com/KillReall666/yaproject/internal/logger"

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

	mylog, err1 := logger2.InitLogger()
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


	r.Post("/update/*", mylog.PostLogger(updateHandler.UpdateMetrics))
	r.Get("/value/*", mylog.GetLogger(getHandler.GetMetrics))

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
