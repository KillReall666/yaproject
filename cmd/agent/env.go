package main

import (
	"github.com/caarlos0/env"
	"log"
)

func setEnv() error {
	cfg := struct {
		Address               string `env:"ADDRESS"`
		DefaultPollInterval   int    `env:"REPORT_INTERVAL"`
		DefaultReportInterval int    `env:"POLL_INTERVAL"`
	}{}

	env.Parse(&cfg)

	if cfg.Address != "" {
		address = cfg.Address
		log.Println("use env variable for address")
	}

	if cfg.DefaultReportInterval != 0 {
		defaultReportInterval = cfg.DefaultReportInterval
		log.Println("use env variable for defaultReportInterval")
	}

	if cfg.DefaultPollInterval != 0 {
		defaultPollInterval = cfg.DefaultPollInterval
		log.Println("use env variable for defaultPollInterval")
	}
	return nil
}
