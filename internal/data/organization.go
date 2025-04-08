package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Organization struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type OrganizationModel struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (m OrganizationModel) Insert(org *Organization) error {
	query := `
		INSERT INTO organizations (name, location)
		VALUES ($1, $2)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, org.Name, org.Location).Scan(&org.ID)
}

func (m OrganizationModel) Get(id int) (*Organization, error) {
	query := `
		SELECT id, name, location
		FROM organizations
		WHERE id = $1
	`
	var org Organization
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&org.ID, &org.Name, &org.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("organization not found: id %d", id)
		}
		return nil, err
	}

	return &org, nil
}

func (m OrganizationModel) Update(org *Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, location = $2
		WHERE id = $3
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, org.Name, org.Location, org.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no organization found with id %d", org.ID)
	}

	return nil
}

func (m OrganizationModel) Delete(id int) error {
	query := `
		DELETE FROM organizations
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
		return fmt.Errorf("no organization found with id %d", id)
	}

	return nil
}

func (m OrganizationModel) List() ([]*Organization, error) {
	query := `
		SELECT id, name, location
		FROM organizations
		ORDER BY name
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Location); err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgs, nil
}
