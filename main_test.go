package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"records/internal/models"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		year INTEGER,
		genre TEXT,
		notes TEXT
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestGetAllRecordsSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec(`INSERT INTO records (title, artist, year, genre, notes) VALUES ('Abbey Road', 'The Beatles', 1969, 'Rock', 'Classic')`)
	if err != nil {
		t.Fatalf("failed to insert test record: %v", err)
	}

	router := createRouter(db)
	req := httptest.NewRequest(http.MethodGet, "/records", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}

	var records []models.Record
	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Title != "Abbey Road" {
		t.Errorf("expected title 'Abbey Road', got '%s'", records[0].Title)
	}
}

func TestGetAllRecordsProblemDetailsOnDBError(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	// Close immediately so queries fail
	db.Close()

	router := createRouter(db)
	req := httptest.NewRequest(http.MethodGet, "/records", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != models.ContentTypeProblemJSON {
		t.Errorf("expected Content-Type '%s', got '%s'", models.ContentTypeProblemJSON, contentType)
	}

	var prob models.ProblemDetails
	if err := json.NewDecoder(rec.Body).Decode(&prob); err != nil {
		t.Fatalf("failed to decode problem details response: %v", err)
	}

	if prob.Status != http.StatusInternalServerError {
		t.Errorf("expected status in body %d, got %d", http.StatusInternalServerError, prob.Status)
	}
	if prob.Title != "Internal Server Error" {
		t.Errorf("expected title 'Internal Server Error', got '%s'", prob.Title)
	}
	if prob.Detail == "" {
		t.Errorf("expected non-empty detail in problem details")
	}
	if prob.Instance != "/records" {
		t.Errorf("expected instance '/records', got '%s'", prob.Instance)
	}
}
