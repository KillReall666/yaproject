package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type RunFileIo struct {
	Interval int    `env:"STORE_INTERVAL"`
	Path     string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
	//ShutdownChan chan os.Signal
}

const (
	defaultInterval = 300
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

	//	cfg.ShutdownChan = make(chan os.Signal, 1)
	//	signal.Notify(cfg.ShutdownChan, syscall.SIGINT, syscall.SIGTERM)

	return cfg
}
