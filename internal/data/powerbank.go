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

type PowerbankStatus string

const (
	PowerbankStatusRented    PowerbankStatus = "rented"
	PowerbankStatusAvailable PowerbankStatus = "available"
	PowerbankStatusCharging  PowerbankStatus = "charging"
)

type Powerbank struct {
	ID               int             `json:"id"`
	CurrentStationID int             `json:"current_station_id"`
	Status           PowerbankStatus `json:"status"`
}

type PowerbankModel struct {
	DB    *sql.DB
	Redis *redis.Client
}

func ValidatePowerbank(v *validator.Validator, p *Powerbank) {
	v.Check(p.CurrentStationID > 0, "current_station_id", "must be a positive integer")
	if !p.Status.IsValid() {
		v.AddError("status", "must be one of: rented, available, charging")
	}
}

func (ps PowerbankStatus) IsValid() bool {
	switch ps {
	case PowerbankStatusRented, PowerbankStatusAvailable, PowerbankStatusCharging:
		return true
	default:
		return false
	}
}

func (m PowerbankModel) ClarifyStatus(id int) (PowerbankStatus, error) {
	query := `
		SELECT status
		FROM powerbanks
		WHERE id = $1;
	`
	var status PowerbankStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("powerbank с id %d не найден", id)
		}
		return "", err
	}

	if !status.IsValid() {
		return "", fmt.Errorf("powerbank с id %d имеет некорректный статус: %s", id, status)
	}

	return status, nil
}

func (m PowerbankModel) Insert(p *Powerbank) error {
	query := `
		INSERT INTO powerbanks (current_station_id, status)
		VALUES ($1, $2)
		RETURNING id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, p.CurrentStationID, p.Status).Scan(&p.ID)
}

func (m PowerbankModel) Get(id int) (*Powerbank, error) {
	query := `
		SELECT id, current_station_id, status
		FROM powerbanks
		WHERE id = $1;
	`
	var p Powerbank
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.CurrentStationID, &p.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("powerbank не найден: id %d", id)
		}
		return nil, err
	}

	return &p, nil
}

func (m PowerbankModel) Update(p *Powerbank) error {
	query := `
		UPDATE powerbanks
		SET current_station_id = $1,
			status = $2
		WHERE id = $3;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, p.CurrentStationID, p.Status, p.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("powerbank с id %d не найден", p.ID)
	}

	return nil
}

func (m PowerbankModel) Delete(id int) error {
	query := `
		DELETE FROM powerbanks
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
		return fmt.Errorf("powerbank с id %d не найден", id)
	}

	return nil
}

func (m PowerbankModel) List() ([]*Powerbank, error) {
	query := `
		SELECT id, current_station_id, status
		FROM powerbanks
		ORDER BY id;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var powerbanks []*Powerbank
	for rows.Next() {
		var p Powerbank
		if err := rows.Scan(&p.ID, &p.CurrentStationID, &p.Status); err != nil {
			return nil, err
		}
		powerbanks = append(powerbanks, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return powerbanks, nil
}
