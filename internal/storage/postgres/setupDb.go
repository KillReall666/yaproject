package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Database struct {
	db *pgxpool.Pool
}

const createTableQuery = `
        CREATE TABLE IF NOT EXISTS metrics(
            name TEXT PRIMARY KEY,
            counter BIGINT,
            gauge DOUBLE PRECISION,
            created_at TIMESTAMP DEFAULT now()
        )
    `

func NewDB(connString string) (*Database, error) {
	if connString == "" {
		return nil, nil
	}

	conn, err := pgxpool.New(context.Background(), connString) //cfg
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	_, err = conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating metrics table: %v", err)
	}

	//TODO: Добавить логгер
	log.Println("Database loaded.")

	return &Database{db: conn}, nil
}

func (d *Database) DBStatusCheck() error {
	err := d.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GaugeSetter(ctx context.Context, name string, gauge float64) error {
	insertQuery := `
                INSERT INTO metrics (name,  gauge)
                VALUES ($1, $2)
                ON CONFLICT (name) 
                DO UPDATE
                SET gauge = EXCLUDED.gauge
            `
  
	_, err := d.db.Exec(ctx, insertQuery, name, gauge)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("error when inserting value to database: %v", pgErr)
		}
	}

	return nil
}

func (d *Database) CountSetter(ctx context.Context, name string, count int64) error {
	insertQuery := `
                INSERT INTO metrics (name, counter)
                VALUES ($1, $2)
                ON CONFLICT (name)
                DO UPDATE
                SET counter = metrics.counter + EXCLUDED.counter
            `
	_, err := d.db.Exec(ctx, insertQuery, name, count)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("error when inserting value to database: %v", pgErr)
		}
	}

	return nil
}

func (d *Database) GaugeGetter(ctx context.Context, key string) (float64, error) {
	var gauge float64

	row := d.db.QueryRow(ctx, "SELECT gauge FROM metrics WHERE name = $1 limit 1", key)

	err := row.Scan(&gauge)
	if err != nil {
		return 0, fmt.Errorf("error get gauge metric: %v", err)
	}

	return gauge, nil
}

func (d *Database) CountGetter(ctx context.Context, key string) (int64, error) {
	var counter int64

	row := d.db.QueryRow(ctx, "SELECT counter FROM metrics WHERE name = $1 limit 1", key)

	err := row.Scan(&counter)
	if err != nil {
		return 0, fmt.Errorf("error get count metric: %v", err)

	}

	return counter, nil
}

func (d *Database) GetAllMetrics() string {
	return ""
}
