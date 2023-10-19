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

func NewDB(connString string) (*Database, *pgx.Conn, error) {
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

func (d *Database) GaugeSetter(name string, gauge float64) error {
	insertQuery := `
                INSERT INTO metrics (name,  gauge)
                VALUES ($1, $2)
                ON CONFLICT (name) 
                DO UPDATE
                SET gauge = EXCLUDED.gauge
            `

	_, err := d.db.Exec(context.Background(), insertQuery, name, gauge)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("error when inserting value to database: %v", pgErr)
		}
	}

	return nil
}

func (d *Database) CountSetter(name string, count int64) error {
	insertQuery := `
                INSERT INTO metrics (name, counter)
                VALUES ($1, $2)
                ON CONFLICT (name)
                DO UPDATE
                SET counter = metrics.counter + EXCLUDED.counter
            `

	_, err := d.db.Exec(context.Background(), insertQuery, name, count)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("error when inserting value to database: %v", pgErr)
		}
	}

	return nil
}

func (d *Database) GaugeGetter(key string) (float64, error) {
	intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	var gauge float64
	found := false
	for _, interval := range intervals {
		rows, err := d.db.Query(context.Background(), "SELECT gauge FROM metrics WHERE name = $1", key)
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

func (d *Database) CountGetter(key string) (int64, error) {
	intervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	var counter int64
	found := false

	for _, interval := range intervals {
		rows, err := d.db.Query(context.Background(), "SELECT counter FROM metrics WHERE name = $1", key)
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
		return 0, errors.New("count value not found in db")
	}

	return counter, nil
}

func (d *Database) GetAllMetrics() string {
	return ""
}
