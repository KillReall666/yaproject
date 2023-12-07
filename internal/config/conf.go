package config

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type RunConfig struct {
	Address               string `env:"ADDRESS"`
	DefaultPollInterval   int    `env:"REPORT_INTERVAL"`
	DefaultReportInterval int    `env:"POLL_INTERVAL"`
	DefaultDBConnStr      string `env:"DATABASE_DSN"`
	Interval              int    `env:"STORE_INTERVAL"`
	Path                  string `env:"FILE_STORAGE_PATH"`
	Restore               bool   `env:"RESTORE"`
	HashKey               string `env:"KEY"`
	RateLimit             int    `env:"RATE_LIMIT"`
	UseDB                 bool
}

const (
	defaultServer             = ":8080"
	defaultPollInterval       = 2
	defaultReportInterval     = 10
	defaultConnStr            = "host=localhost user=Mr8 password=Rammstein12! dbname=yaproject_db sslmode=disable"
	defaultSaveOnDiskInterval = 300
	defaultPathOfFile         = "./metrics-postgres.json"
	defaultRestore            = true
)

func LoadForAgent() RunConfig {
	cfg := RunConfig{}

	flag.IntVar(&cfg.DefaultPollInterval, "p", defaultPollInterval, "metrics update interval in seconds")
	flag.IntVar(&cfg.DefaultReportInterval, "r", defaultReportInterval, "metrics sending interval in seconds")
	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", 5, "pool workers limit")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}

func LoadForServer() (RunConfig, error) {
	cfg := RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.StringVar(&cfg.DefaultDBConnStr, "d", defaultConnStr, "connection string")
	flag.IntVar(&cfg.Interval, "i", defaultSaveOnDiskInterval, "time interval in seconds after which the current server readings are saved to disk")
	flag.StringVar(&cfg.Path, "f", defaultPathOfFile, "full name of the file where the current values are saved")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore, "load or not previously saved values from specified files when starting the server")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	cfg.UseDB = true
	if cfg.DefaultDBConnStr == "" {
		cfg.UseDB = false
		err = errors.New("metric storage switched to memory, the database is not connected")
		return cfg, err
	}

	return cfg, err
}
