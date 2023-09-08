package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KillReall666/yaproject/internal/handlers/update"
	"github.com/KillReall666/yaproject/internal/service"
	"github.com/KillReall666/yaproject/internal/storage"
)

const (
	port     = ":8080"
	endpoint = "/update/"
)

func main() {
	server := http.NewServeMux()

	store := storage.NewMemStorage()

	serv := service.NewService(store)

	updateHandler := update.NewHandler(serv)
	server.HandleFunc(endpoint, updateHandler.PostHandle)

	log.Printf("Starting http server to serve metrics at port %s", port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Printf("server is down: %v", err)
		panic(fmt.Errorf("server is down: %v", err))
	}

}
