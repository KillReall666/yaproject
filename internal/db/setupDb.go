package repo

import (
	"context"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"log"
)

func GetDB() (*pgx.Conn, error) {
	connString := config.LoadDbConfig()
	cfg, err := pgx.ParseConfig(connString.DefaultConnStr)
	if err != nil {
		log.Println("ошибка при разборе строки подключения:", err)
		return nil, err
	}
	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Println("ошибка при подключении к БД:", err)
		return nil, err
	}
	defer conn.Close(context.Background())
	log.Println("connected to db!")
	return conn, nil

}
