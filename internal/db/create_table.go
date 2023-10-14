package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

const createTableQuery = `
        CREATE TABLE IF NOT EXISTS metrics(
            name TEXT PRIMARY KEY,
            counter BIGINT,
            gauge DOUBLE PRECISION,
            created_at TIMESTAMP DEFAULT now()
        )
    `

func (d *Database) CreateMetricsTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		return fmt.Errorf("Error creating metrics table: %v", err)
	}
	return nil
}
