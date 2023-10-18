package db

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"time"
)

func (d *Database) GaugeSetter(name string, gauge float64) error {
	err := retry.Do(
		func() error {
			insertQuery := `
                INSERT INTO metrics (name,  gauge)
                VALUES ($1, $2)
                ON CONFLICT (name) 
                DO UPDATE
                SET gauge = EXCLUDED.gauge
            `

			_, err := d.db.Exec(context.Background(), insertQuery, name, gauge)
			if err != nil {
				return fmt.Errorf("error inserting gauge metric: %v", err)
			}

			return nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)

	return err
}

func (d *Database) CountSetter(name string, count int64) error {
	err := retry.Do(
		func() error {
			insertQuery := `
                INSERT INTO metrics (name, counter)
                VALUES ($1, $2)
                ON CONFLICT (name)
                DO UPDATE
                SET counter = metrics.counter + EXCLUDED.counter
            `

			_, err := d.db.Exec(context.Background(), insertQuery, name, count)
			if err != nil {
				return fmt.Errorf("error inserting counter metric: %v", err)
			}

			return nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)

	return err
}
