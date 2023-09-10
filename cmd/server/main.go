package main

import (
	"fmt"
	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const (
	port           = ":8080"
	endpointUpdate = "/update/*"
	endpointCount  = "/value/*"
	endpointHtml   = "/"
)

func main() {
	router := chi.NewRouter()
	store := storage.NewMemStorage()
	serv := service.NewService(store)
	updateHandler := update.NewHandler(serv)

	router.Post(endpointUpdate, updateHandler.PostHandle)
	router.Get(endpointCount, updateHandler.GetHandle)
	router.HandleFunc(endpointHtml, updateHandler.HtmlHandle)
	log.Printf("Starting http server to serve metrics at port%s ", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Printf("server is down: %v", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
