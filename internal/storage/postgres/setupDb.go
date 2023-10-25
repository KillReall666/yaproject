package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type Database struct {
	db *pgx.Conn
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
	cfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %v", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	_, err = conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating metrics table: %v", err)
	}

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
	intervals := []time.Duration{1 * time.Second, 1 * time.Second, 3 * time.Second, 5 * time.Second}
	var gauge float64
	found := false
	for _, interval := range intervals {
		rows, err := d.db.Query(ctx, "SELECT gauge FROM metrics WHERE name = $1", key)
		if err != nil {
			return 0, err
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&gauge)
			if err != nil {
				return 0, fmt.Errorf("error get gauge metric: %v", err)
			}
			found = true
		}

		if found {
			break
		}

		time.Sleep(interval)
	}

	if !found {
		return 0, errors.New("gauge value not found in postgres")
	}

	return gauge, nil
}

func (d *Database) CountGetter(ctx context.Context, key string) (int64, error) {
	intervals := []time.Duration{1 * time.Second, 1 * time.Second, 3 * time.Second, 5 * time.Second}

	var counter int64
	found := false

	for _, interval := range intervals {
		rows, err := d.db.Query(ctx, "SELECT counter FROM metrics WHERE name = $1", key)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&counter)
			if err != nil {
				return 0, fmt.Errorf("error get count metric: %v", err)
			}
			found = true
		}

		if found {
			break
		}

		time.Sleep(interval)
	}

	if !found {
		return 0, errors.New("count value not found in postgres")
	}

	return counter, nil
}

func (d *Database) GetAllMetrics() string {
	return ""
}
