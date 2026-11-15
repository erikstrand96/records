package services

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"records/internal/models"

	_ "modernc.org/sqlite"
)

// setupTestDB creates a temporary SQLite database with the records schema for testing.
func setupTestDB(t *testing.T) (*RecordService, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}

	// Create the records table manually instead of relying on migrations.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			artist TEXT NOT NULL,
			year INTEGER,
			genre TEXT,
			notes TEXT
		)
	`)
	if err != nil {
		db.Close()
		os.Remove(dbPath)
		t.Fatalf("creating records table: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return NewRecordService(db), cleanup
}

func TestRecordService_CreateAndGetAll(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	year := 2023
	genre := "Rock"
	notes := "First pressing"

	req := models.CreateRecordRequest{
		Title:  "Dark Side of the Moon",
		Artist: "Pink Floyd",
		Year:   &year,
		Genre:  &genre,
		Notes:  &notes,
	}

	created, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created record should have a non-zero ID")
	}
	if created.Title != req.Title {
		t.Errorf("title = %q, want %q", created.Title, req.Title)
	}

	records, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("GetAll returned %d records, want 1", len(records))
	}
	if records[0].Title != req.Title {
		t.Errorf("GetAll[0].Title = %q, want %q", records[0].Title, req.Title)
	}
}

func TestRecordService_GetByID(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	year := 1991
	req := models.CreateRecordRequest{
		Title:  "Nevermind",
		Artist: "Nirvana",
		Year:   &year,
	}

	created, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	found, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if found.Title != req.Title {
		t.Errorf("Title = %q, want %q", found.Title, req.Title)
	}
}

func TestRecordService_GetByID_NotFound(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := svc.GetByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent record, got nil")
	}
}

func TestRecordService_Update(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	year := 2000
	req := models.CreateRecordRequest{
		Title:  "Original Title",
		Artist: "Artist",
		Year:   &year,
	}

	created, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	newTitle := "Updated Title"
	affected, err := svc.Update(created.ID, models.UpdateRecordRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if affected != 1 {
		t.Errorf("Update affected %d rows, want 1", affected)
	}

	updated, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("Title = %q, want %q", updated.Title, newTitle)
	}
}

func TestRecordService_Update_NotFound(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	newTitle := "Nope"
	affected, err := svc.Update(999, models.UpdateRecordRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if affected != 0 {
		t.Errorf("Update affected %d rows, want 0", affected)
	}
}

func TestRecordService_Update_NoFields(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := svc.Create(models.CreateRecordRequest{
		Title:  "Test",
		Artist: "Test Artist",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	_, err = svc.Update(1, models.UpdateRecordRequest{})
	if err == nil {
		t.Fatal("expected error when no fields are provided, got nil")
	}
}

func TestRecordService_Delete(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	req := models.CreateRecordRequest{
		Title:  "To Be Deleted",
		Artist: "Artist",
	}

	created, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	affected, err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if affected != 1 {
		t.Errorf("Delete affected %d rows, want 1", affected)
	}

	_, err = svc.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestRecordService_Delete_NotFound(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	affected, err := svc.Delete(999)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if affected != 0 {
		t.Errorf("Delete affected %d rows, want 0", affected)
	}
}

func TestRecordService_GetAll_Empty(t *testing.T) {
	svc, cleanup := setupTestDB(t)
	defer cleanup()

	records, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if records == nil {
		t.Error("GetAll should return empty slice, not nil")
	}
	if len(records) != 0 {
		t.Errorf("GetAll returned %d records, want 0", len(records))
	}
}
