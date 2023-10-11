package db

import (
	"context"
	"github.com/KillReall666/yaproject/internal/config"
	"github.com/jackc/pgx/v5"
	"log"
)

type Database struct {
	db *pgx.Conn
}

func GetDB() (*Database, error) {
	connString := config.LoadDBConfig()

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

	log.Println("подключение с БД установлено!")

	return &Database{db: conn}, nil
}

func (d *Database) DBStatusCheck() error {
	err := d.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
