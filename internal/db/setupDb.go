package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	db *pgx.Conn
}

func GetDB(connString string) (*Database, *pgx.Conn, error) {
	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing connection string: %v", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &Database{db: conn}, conn, nil
}

func (d *Database) DBStatusCheck() error {
	err := d.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
