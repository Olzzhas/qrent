package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/olzzhas/qrent/pkg/validator"
	"time"
)

type Station struct {
	ID    int `json:"id"`
	OrgID int `json:"org_id"`
}

type StationModel struct {
	DB    *sql.DB
	Redis *redis.Client
}

func ValidateStation(v *validator.Validator, s *Station) {
	v.Check(s.OrgID > 0, "org_id", "must be a positive integer")
}

func (m StationModel) Insert(station *Station) error {
	query := `
		INSERT INTO stations (org_id)
		VALUES ($1)
		RETURNING id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, station.OrgID).Scan(&station.ID)
}

func (m StationModel) Get(id int) (*Station, error) {
	query := `
		SELECT id, org_id
		FROM stations
		WHERE id = $1;
	`
	var station Station
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&station.ID, &station.OrgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("station not found: id %d", id)
		}
		return nil, err
	}

	return &station, nil
}

func (m StationModel) Update(station *Station) error {
	query := `
		UPDATE stations
		SET org_id = $1
		WHERE id = $2;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, station.OrgID, station.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no station found with id %d", station.ID)
	}

	return nil
}

func (m StationModel) Delete(id int) error {
	query := `
		DELETE FROM stations
		WHERE id = $1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no station found with id %d", id)
	}

	return nil
}

func (m StationModel) List() ([]*Station, error) {
	query := `
		SELECT id, org_id
		FROM stations
		ORDER BY id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []*Station
	for rows.Next() {
		var station Station
		if err := rows.Scan(&station.ID, &station.OrgID); err != nil {
			return nil, err
		}
		stations = append(stations, &station)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
