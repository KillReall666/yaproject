package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type RunConfig struct {
	Address               string `env:"ADDRESS"`
	DefaultPollInterval   int    `env:"REPORT_INTERVAL"`
	DefaultReportInterval int    `env:"POLL_INTERVAL"`
}

const (
	defaultServer         = ":8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

func LoadAgentConfig() RunConfig {
	cfg := RunConfig{}

	flag.IntVar(&cfg.DefaultPollInterval, "p", defaultPollInterval, "metrics html interval in seconds")
	flag.IntVar(&cfg.DefaultReportInterval, "r", defaultReportInterval, "metrics sending interval in seconds")
	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}

func LoadServerConfig() RunConfig {
	cfg := RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}
