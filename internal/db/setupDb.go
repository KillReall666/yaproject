package db

import (
	"context"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	db *pgx.Conn
}

func GetDB(l *logger.Logger, connString string) (*Database, *pgx.Conn, error) {
	if connString == "" {
		l.LogInfo("база данных подключена не будет. хранение метрик будет произведено в памяти.")
		return nil, nil, nil
	}
	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		l.LogInfo("ошибка при разборе строки подключения:", err)
		//return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		l.LogInfo("ошибка при подключении к БД:", err)
		//return nil, err
	}

	l.LogInfo("подключение с БД установлено!")

	return &Database{db: conn}, conn, nil
}

func (d *Database) DBStatusCheck() error {
	err := d.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
