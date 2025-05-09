package data

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
)

var (
	ErrInvalidForeignKey = errors.New("invalid foreign key: related record does not exist")
	ErrRecordNotFound    = errors.New("record not found")
)

type Models struct {
	Organization OrganizationModel
	Powerbank    PowerbankModel
	Station      StationModel
}

func NewModels(db *sql.DB, redis *redis.Client) Models {
	return Models{
		Organization: OrganizationModel{DB: db, Redis: redis},
		Powerbank:    PowerbankModel{DB: db, Redis: redis},
		Station:      StationModel{DB: db, Redis: redis},
	}
}
