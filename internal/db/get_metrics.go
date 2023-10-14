package db

import (
	"context"
	"fmt"
)

func (d *Database) GetGaugeMetricsFromDB(name string) (float64, error) {
	rows, err := d.db.Query(context.Background(), "SELECT gauge FROM metrics WHERE name = $1", name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var gauge float64

		err := rows.Scan(&gauge)
		if err != nil {
			return gauge, nil
		}
	}
	return 0, nil
}

func (d *Database) GetCounterMetricsFromDB(name string) (int64, error) {
	rows, err := d.db.Query(context.Background(), "SELECT counter FROM metrics WHERE name = $1", name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var counter int64

		err := rows.Scan(&counter)
		if err != nil {
			fmt.Println("Error reading row:", err)
			return counter, nil
		}
	}
	return 0, nil
}
