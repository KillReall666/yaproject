package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type RunFileIo struct {
	Interval int    `env:"STORE_INTERVAL"`
	Path     string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
}

const (
	defaultInterval = 15
	defaultPath     = "./metrics-db.json"
	defaultRestore  = true
)

func LoadFileIoConf() RunFileIo {
	cfg := RunFileIo{}
	flag.IntVar(&cfg.Interval, "i", defaultInterval, "time interval in seconds after which the current server readings are saved to disk")
	flag.StringVar(&cfg.Path, "f", defaultPath, "full name of the file where the current values are saved")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore, "load or not previously saved values from specified files when starting the server")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	return cfg
}
