package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type DBConf struct {
	DefaultConnStr string `env:"DATABASE_DSN"`
}

const defaultConnStr = "host=localhost user=Mr8 password=Rammstein12! dbname=yaproject_db sslmode=disable"

func LoadDBConfig() DBConf {
	str := DBConf{}
	flag.StringVar(&str.DefaultConnStr, "d", defaultConnStr, "connection string")
	flag.Parse()

	err := env.Parse(&str)
	if err != nil {
		log.Println("ошибка при парсинге переменной: ", err)
	}

	return str
}
