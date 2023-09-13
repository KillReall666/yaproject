package config

import (
	"flag"
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/caarlos0/env"
	"log"
)

const (
	defaultServer         = ":8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

func LoadAgentConfig() model.RunConfig {
	cfg := model.RunConfig{}

	flag.IntVar(&cfg.DefaultPollInterval, "p", defaultPollInterval, "metrics update interval in seconds")
	flag.IntVar(&cfg.DefaultReportInterval, "r", defaultReportInterval, "metrics sending interval in seconds")
	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}

func LoadServerConfig() model.RunConfig {
	cfg := model.RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}
