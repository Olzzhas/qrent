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
	ID        int       `json:"id"`
	OrgID     int       `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
        RETURNING id, created_at, updated_at
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, station.OrgID).
		Scan(&station.ID, &station.CreatedAt, &station.UpdatedAt)
}

func (m StationModel) Get(id int) (*Station, error) {
	query := `
        SELECT id, org_id, created_at, updated_at
        FROM stations
        WHERE id = $1
    `
	var station Station
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).
		Scan(&station.ID, &station.OrgID, &station.CreatedAt, &station.UpdatedAt)
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
        SET org_id = $1,
            updated_at = NOW() -- или полагайтесь на триггер, если он есть
        WHERE id = $2
        RETURNING updated_at
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, station.OrgID, station.ID).
		Scan(&station.UpdatedAt)
	if err != nil {
		// Если строка не найдена, будет sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no station found with id %d", station.ID)
		}
		return err
	}

	return nil
}

func (m StationModel) Delete(id int) error {
	query := `
        DELETE FROM stations
        WHERE id = $1
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
        SELECT id, org_id, created_at, updated_at
        FROM stations
        ORDER BY id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stations := make([]*Station, 0)
	for rows.Next() {
		var station Station
		if err := rows.Scan(
			&station.ID,
			&station.OrgID,
			&station.CreatedAt,
			&station.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stations = append(stations, &station)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
