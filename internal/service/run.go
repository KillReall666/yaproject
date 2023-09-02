package service

import (
	"github.com/KillReall666/yaproject/internal/model"
	"log"
	"net/http"
)

var (
	port     = ":8080"
	endpoint = "/update/"
)

var storage *model.MetricsStorage

func Run() error {
	log.Printf("Starting http server to serve metrics at port %s", port)
	server := http.NewServeMux()
	server.HandleFunc(endpoint, PostHandler)
	storage = newMetricStorage()
	return http.ListenAndServe(port, server)

}
