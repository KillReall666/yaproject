package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"time"
)

func (d Database) GaugeGetter(key string) (float64, error) {
	var gauge float64
	err := retry.Do(
		func() error {
			rows, err := d.db.Query(context.Background(), "SELECT gauge FROM metrics WHERE name = $1", key)
			if err != nil {
				return err
			}
			defer rows.Close()

			found := false

			for rows.Next() {
				err = rows.Scan(&gauge)
				if err != nil {
					return fmt.Errorf("error get gauge metric: %v", err)
				}
				found = true
			}
			if !found {
				return errors.New("gauge value not found in db")
			}

			return nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)

	return gauge, err
}

func (d *Database) CountGetter(key string) (int64, error) {
	var counter int64
	err := retry.Do(
		func() error {
			rows, err := d.db.Query(context.Background(), "SELECT counter FROM metrics WHERE name = $1", key)
			if err != nil {
				return err
			}
			defer rows.Close()

			found := false

			for rows.Next() {
				err = rows.Scan(&counter)
				if err != nil {
					return fmt.Errorf("error get count metric: %v", err)
				}
				found = true
			}

			if !found {
				return errors.New("count value not found in db")
			}

			return nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)
	return counter, err
}

func (d *Database) GetAllMetrics() string {
	return ""
}
