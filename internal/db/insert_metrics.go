package db

import (
	"context"
	"fmt"
)

func (d *Database) InsertGaugeMetrics(name string, gauge *float64) error {
	insertQuery := `
        INSERT INTO metrics (name,  gauge)
        VALUES ($1, $2)
		ON CONFLICT (name) 
		DO UPDATE
		SET gauge = EXCLUDED.gauge
    `

	_, err := d.db.Exec(context.Background(), insertQuery, name, gauge)
	if err != nil {
		return fmt.Errorf("error inserting  gauge metric: %v", err)
	}

	return nil
}

func (d *Database) InsertCountMetrics(name string, counter *int64) error {
	insertQuery := `
        INSERT  INTO metrics (name, counter)
        VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE
		SET counter = EXCLUDED.counter
    `

	_, err := d.db.Exec(context.Background(), insertQuery, name, counter)
	if err != nil {
		return fmt.Errorf("error inserting counter metric: %v", err)
	}

	return nil
}
