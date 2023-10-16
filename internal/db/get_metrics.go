package db

import (
	"context"
	"log"
)

func (d *Database) GetGaugeFromDB(name string) (float64, error) {
	rows, err := d.db.Query(context.Background(), "SELECT gauge FROM metrics WHERE name = $1", name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var gauge float64

	for rows.Next() {
		err := rows.Scan(&gauge)
		if err != nil {
			return 0, nil
		}
	}
	return gauge, nil
}

func (d *Database) GetCounterFromDB(name string) (int64, error) {
	rows, err := d.db.Query(context.Background(), "SELECT counter FROM metrics WHERE name = $1", name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var counter int64

	for rows.Next() {
		err := rows.Scan(&counter)
		if err != nil {
			log.Println("Error reading row:", err)
			return 0, nil
		}
	}
	return counter, nil
}
