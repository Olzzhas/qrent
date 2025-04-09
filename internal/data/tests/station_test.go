// station_model_test.go
package data_test

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/olzzhas/qrent/internal/data"
)

// Test Insert успеха
func TestStationModel_Insert_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}

	station := &data.Station{
		OrgID: 10,
	}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO stations (org_id)
        VALUES ($1)
        RETURNING id, created_at, updated_at
    `)).
		WithArgs(station.OrgID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(5, now, now))

	if err = model.Insert(station); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if station.ID != 5 {
		t.Errorf("expected station.ID=5, got %d", station.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Test Insert при ошибке внешнего ключа
func TestStationModel_Insert_InvalidForeignKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	station := &data.Station{OrgID: 99999} // Предполагаем неверный OrgID

	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO stations (org_id)
        VALUES ($1)
        RETURNING id, created_at, updated_at
    `)).
		WithArgs(station.OrgID).
		WillReturnError(&pq.Error{Code: "23503"})

	err = model.Insert(station)
	if err == nil {
		t.Errorf("expected error due to invalid foreign key, got nil")
	}
	expectedErr := data.ErrInvalidForeignKey
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

// Test Get успеха
func TestStationModel_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	now := time.Now()

	// Ожидаем запрос, возвращающий корректную запись
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, org_id, created_at, updated_at
        FROM stations
        WHERE id = $1
    `)).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "org_id", "created_at", "updated_at"}).
			AddRow(5, 10, now, now))

	station, err := model.Get(5)
	if err != nil {
		t.Errorf("unexpected error in Get: %s", err)
	}
	if station.ID != 5 || station.OrgID != 10 {
		t.Errorf("unexpected station data: %+v", station)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Test Get: не найден
func TestStationModel_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, org_id, created_at, updated_at
        FROM stations
        WHERE id = $1
    `)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = model.Get(999)
	if err == nil {
		t.Errorf("expected error for non-existent station, got nil")
	}
	expected := fmt.Sprintf("station not found: id %d", 999)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

// Test Update успеха
func TestStationModel_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	station := &data.Station{
		ID:    5,
		OrgID: 20,
	}

	newTime := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
        UPDATE stations
        SET org_id = $1,
            updated_at = NOW()
        WHERE id = $2
        RETURNING updated_at
    `)).
		WithArgs(station.OrgID, station.ID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(newTime))

	err = model.Update(station)
	if err != nil {
		t.Errorf("unexpected error in Update: %s", err)
	}
	if !station.UpdatedAt.Equal(newTime) {
		t.Errorf("expected updated_at %v, got %v", newTime, station.UpdatedAt)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Test Update: не найден
func TestStationModel_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	station := &data.Station{
		ID:    999,
		OrgID: 20,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
        UPDATE stations
        SET org_id = $1,
            updated_at = NOW()
        WHERE id = $2
        RETURNING updated_at
    `)).
		WithArgs(station.OrgID, station.ID).
		WillReturnError(sql.ErrNoRows)

	err = model.Update(station)
	if err == nil {
		t.Errorf("expected error for non-existent station in Update, got nil")
	}
	expected := fmt.Sprintf("no station found with id %d", station.ID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

// Test Update: invalid foreign key
func TestStationModel_Update_InvalidForeignKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	station := &data.Station{
		ID:    5,
		OrgID: 99999, // несуществующий OrgID
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
        UPDATE stations
        SET org_id = $1,
            updated_at = NOW()
        WHERE id = $2
        RETURNING updated_at
    `)).
		WithArgs(station.OrgID, station.ID).
		WillReturnError(&pq.Error{Code: "23503"})

	err = model.Update(station)
	if err == nil {
		t.Errorf("expected error due to invalid foreign key in Update, got nil")
	}
	expectedErr := data.ErrInvalidForeignKey
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

// Test Delete успеха
func TestStationModel_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	stationID := 5

	mock.ExpectExec(regexp.QuoteMeta(`
        DELETE FROM stations
        WHERE id = $1
    `)).
		WithArgs(stationID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = model.Delete(stationID)
	if err != nil {
		t.Errorf("unexpected error in Delete: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Test Delete: не найден
func TestStationModel_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	stationID := 999

	mock.ExpectExec(regexp.QuoteMeta(`
        DELETE FROM stations
        WHERE id = $1
    `)).
		WithArgs(stationID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 строк затронуто

	err = model.Delete(stationID)
	if err == nil {
		t.Errorf("expected error for non-existent station in Delete, got nil")
	}
	expected := fmt.Sprintf("no station found with id %d", stationID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

// Test List успеха
func TestStationModel_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания sqlmock: %s", err)
	}
	defer db.Close()

	model := data.StationModel{DB: db}
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, org_id, created_at, updated_at
        FROM stations
        ORDER BY id
    `)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "org_id", "created_at", "updated_at"}).
			AddRow(1, 10, now, now).
			AddRow(2, 20, now, now))

	stations, err := model.List()
	if err != nil {
		t.Errorf("unexpected error in List: %s", err)
	}
	if len(stations) != 2 {
		t.Errorf("expected 2 stations, got %d", len(stations))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
