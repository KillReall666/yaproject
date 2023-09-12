package main

import (
	"github.com/caarlos0/env"
	"log"
)

func setEnv() {
	cfg := struct {
		Addr string `env:"ADDRESS"`
	}{}

	env.Parse(&cfg)
	if cfg.Addr != "" {
		addr = cfg.Addr
		log.Println("use env variable!")
		return
	}
	log.Println("env not found! use default flag")
}
