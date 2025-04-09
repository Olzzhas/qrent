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

func TestPowerbankStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   data.PowerbankStatus
		expected bool
	}{
		{data.PowerbankStatusRented, true},
		{data.PowerbankStatusAvailable, true},
		{data.PowerbankStatusCharging, true},
		{"unknown", false},
	}

	for _, tt := range tests {
		if got := tt.status.IsValid(); got != tt.expected {
			t.Errorf("PowerbankStatus(%q).IsValid() = %v, want %v", tt.status, got, tt.expected)
		}
	}
}

func TestPowerbankModel_Insert_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}

	p := &data.Powerbank{
		CurrentStationID: 10,
		Status:           data.PowerbankStatusAvailable,
	}

	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO powerbanks (current_station_id, status)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`)).
		WithArgs(p.CurrentStationID, p.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(5, now, now))

	if err = model.Insert(p); err != nil {
		t.Errorf("unexpected error in Insert: %s", err)
	}
	if p.ID != 5 {
		t.Errorf("expected powerbank ID=5, got %d", p.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_Insert_InvalidForeignKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	p := &data.Powerbank{CurrentStationID: 99999, Status: data.PowerbankStatusRented}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO powerbanks (current_station_id, status)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`)).
		WithArgs(p.CurrentStationID, p.Status).
		WillReturnError(&pq.Error{Code: "23503"})

	err = model.Insert(p)
	if err == nil {
		t.Errorf("expected error due to invalid foreign key, got nil")
	}
	expectedErr := data.ErrInvalidForeignKey
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestPowerbankModel_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, current_station_id, status, created_at, updated_at
		FROM powerbanks
		WHERE id = $1
	`)).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "current_station_id", "status", "created_at", "updated_at"}).
			AddRow(5, 10, data.PowerbankStatusCharging, now, now))

	p, err := model.Get(5)
	if err != nil {
		t.Errorf("unexpected error in Get: %s", err)
	}
	if p.ID != 5 || p.CurrentStationID != 10 || p.Status != data.PowerbankStatusCharging {
		t.Errorf("unexpected powerbank data: %+v", p)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, current_station_id, status, created_at, updated_at
		FROM powerbanks
		WHERE id = $1
	`)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = model.Get(999)
	if err == nil {
		t.Error("expected error for non-existent powerbank in Get, got nil")
	}
	expected := fmt.Sprintf("powerbank не найден: id %d", 999)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestPowerbankModel_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	p := &data.Powerbank{
		ID:               5,
		CurrentStationID: 20,
		Status:           data.PowerbankStatusRented,
	}

	newTime := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE powerbanks
		SET current_station_id = $1,
			status = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`)).
		WithArgs(p.CurrentStationID, p.Status, p.ID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(newTime))

	err = model.Update(p)
	if err != nil {
		t.Errorf("unexpected error in Update: %s", err)
	}
	if !p.UpdatedAt.Equal(newTime) {
		t.Errorf("expected updated_at %v, got %v", newTime, p.UpdatedAt)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	p := &data.Powerbank{
		ID:               999,
		CurrentStationID: 20,
		Status:           data.PowerbankStatusAvailable,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE powerbanks
		SET current_station_id = $1,
			status = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`)).
		WithArgs(p.CurrentStationID, p.Status, p.ID).
		WillReturnError(sql.ErrNoRows)

	err = model.Update(p)
	if err == nil {
		t.Error("expected error for non-existent powerbank in Update, got nil")
	}
	expected := fmt.Sprintf("powerbank with id %d not found", p.ID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestPowerbankModel_Update_InvalidForeignKey(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	p := &data.Powerbank{
		ID:               5,
		CurrentStationID: 99999, // некорректный FK
		Status:           data.PowerbankStatusRented,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE powerbanks
		SET current_station_id = $1,
			status = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`)).
		WithArgs(p.CurrentStationID, p.Status, p.ID).
		WillReturnError(&pq.Error{Code: "23503"})

	err = model.Update(p)
	if err == nil {
		t.Error("expected error due to invalid foreign key in Update, got nil")
	}
	expectedErr := data.ErrInvalidForeignKey
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestPowerbankModel_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	powerbankID := 5

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM powerbanks
		WHERE id = $1
	`)).
		WithArgs(powerbankID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 строка затронута

	err = model.Delete(powerbankID)
	if err != nil {
		t.Errorf("unexpected error in Delete: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	powerbankID := 999

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE FROM powerbanks
		WHERE id = $1
	`)).
		WithArgs(powerbankID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 строк затронуто

	err = model.Delete(powerbankID)
	if err == nil {
		t.Error("expected error for non-existent powerbank in Delete, got nil")
	}
	expected := fmt.Sprintf("powerbank с id %d не найден", powerbankID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestPowerbankModel_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, current_station_id, status, created_at, updated_at
		FROM powerbanks
		ORDER BY id
	`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "current_station_id", "status", "created_at", "updated_at"}).
			AddRow(1, 10, data.PowerbankStatusAvailable, now, now).
			AddRow(2, 20, data.PowerbankStatusRented, now, now))

	powerbanks, err := model.List()
	if err != nil {
		t.Errorf("unexpected error in List: %s", err)
	}
	if len(powerbanks) != 2 {
		t.Errorf("expected 2 powerbanks, got %d", len(powerbanks))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_ClarifyStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT status FROM powerbanks WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).
			AddRow(data.PowerbankStatusCharging))

	status, err := model.ClarifyStatus(1)
	if err != nil {
		t.Errorf("unexpected error in ClarifyStatus: %s", err)
	}
	if status != data.PowerbankStatusCharging {
		t.Errorf("expected status %q, got %q", data.PowerbankStatusCharging, status)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestPowerbankModel_ClarifyStatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating sqlmock: %s", err)
	}
	defer db.Close()

	model := data.PowerbankModel{DB: db}
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT status
		FROM powerbanks
		WHERE id = $1
	`)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = model.ClarifyStatus(999)
	if err == nil {
		t.Error("expected error for non-existent powerbank in ClarifyStatus, got nil")
	}
	expected := data.ErrRecordNotFound
	if !errors.Is(expected, err) {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
