// organization_model_test.go
package data_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/olzzhas/qrent/internal/data"
)

func TestOrganizationModel_Insert_Success(t *testing.T) {
	// Создаем фиктивное подключение к базе данных с помощью sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}

	org := &data.Organization{
		Name:     "Test Org",
		Location: "Test Location",
	}

	now := time.Now()
	// Ожидаем запрос INSERT с заданными аргументами и возвращаем результат
	mock.ExpectQuery(regexp.QuoteMeta(`
        INSERT INTO organizations (name, location)
        VALUES ($1, $2)
        RETURNING id, created_at, updated_at
    `)).
		WithArgs(org.Name, org.Location).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, now, now))

	err = model.Insert(org)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if org.ID != 1 {
		t.Errorf("expected org.ID=1, got %d", org.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestOrganizationModel_Get_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, name, location, created_at, updated_at
        FROM organizations
        WHERE id = $1
    `)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "created_at", "updated_at"}).
			AddRow(1, "Test Org", "Test Location", now, now))

	org, err := model.Get(1)
	if err != nil {
		t.Errorf("unexpected error in Get: %s", err)
	}

	if org.ID != 1 || org.Name != "Test Org" || org.Location != "Test Location" {
		t.Errorf("unexpected organization data: %+v", org)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestOrganizationModel_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, name, location, created_at, updated_at
        FROM organizations
        WHERE id = $1
    `)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = model.Get(999)
	if err == nil {
		t.Error("expected error for non-existent organization, got nil")
	}
	expected := fmt.Sprintf("organization not found: id %d", 999)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestOrganizationModel_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	org := &data.Organization{
		ID:       1,
		Name:     "Old Name",
		Location: "Old Location",
	}

	// Новое время обновления
	newTime := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`
        UPDATE organizations
        SET name = $1,
            location = $2,
            updated_at = NOW()
        WHERE id = $3
        RETURNING updated_at
    `)).
		WithArgs("New Name", "New Location", org.ID).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(newTime))

	// Обновляем данные
	org.Name = "New Name"
	org.Location = "New Location"

	err = model.Update(org)
	if err != nil {
		t.Errorf("unexpected error in Update: %s", err)
	}

	if !org.UpdatedAt.Equal(newTime) {
		t.Errorf("expected updated_at %v, got %v", newTime, org.UpdatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestOrganizationModel_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	org := &data.Organization{
		ID:       999,
		Name:     "Name",
		Location: "Location",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
        UPDATE organizations
        SET name = $1,
            location = $2,
            updated_at = NOW()
        WHERE id = $3
        RETURNING updated_at
    `)).
		WithArgs(org.Name, org.Location, org.ID).
		WillReturnError(sql.ErrNoRows)

	err = model.Update(org)
	if err == nil {
		t.Error("expected error for non-existent organization in Update, got nil")
	}
	expected := fmt.Sprintf("no organization found with id %d", org.ID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestOrganizationModel_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	orgID := 1

	mock.ExpectExec(regexp.QuoteMeta(`
        DELETE FROM organizations
        WHERE id = $1
    `)).
		WithArgs(orgID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = model.Delete(orgID)
	if err != nil {
		t.Errorf("unexpected error in Delete: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestOrganizationModel_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	orgID := 999

	mock.ExpectExec(regexp.QuoteMeta(`
        DELETE FROM organizations
        WHERE id = $1
    `)).
		WithArgs(orgID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 строк затронуто

	err = model.Delete(orgID)
	if err == nil {
		t.Error("expected error for non-existent organization in Delete, got nil")
	}
	expected := fmt.Sprintf("no organization found with id %d", orgID)
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestOrganizationModel_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	model := data.OrganizationModel{DB: db}
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, name, location, created_at, updated_at
        FROM organizations
        ORDER BY name
    `)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "created_at", "updated_at"}).
			AddRow(1, "Org1", "Location1", now, now).
			AddRow(2, "Org2", "Location2", now, now))

	orgs, err := model.List()
	if err != nil {
		t.Errorf("unexpected error in List: %s", err)
	}
	if len(orgs) != 2 {
		t.Errorf("expected 2 organizations, got %d", len(orgs))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
