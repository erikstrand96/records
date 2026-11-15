package services

import (
	"database/sql"
	"errors"
	"fmt"
	"records/internal/models"
)

// RecordService provides database access for Record entities.
type RecordService struct {
	db *sql.DB
}

// NewRecordService creates a new RecordService backed by the given *sql.DB.
func NewRecordService(db *sql.DB) *RecordService {
	return &RecordService{db: db}
}

// GetAll returns all records from the database.
func (s *RecordService) GetAll() ([]models.Record, error) {
	rows, err := s.db.Query("SELECT id, title, artist, year, genre, notes FROM records")
	if err != nil {
		return nil, fmt.Errorf("querying records: %w", err)
	}
	defer rows.Close()

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		if err := rows.Scan(&rec.ID, &rec.Title, &rec.Artist, &rec.Year, &rec.Genre, &rec.Notes); err != nil {
			return nil, fmt.Errorf("scanning record: %w", err)
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating records: %w", err)
	}

	if records == nil {
		records = []models.Record{}
	}

	return records, nil
}

// GetByID returns a single record by its ID.
func (s *RecordService) GetByID(id int64) (*models.Record, error) {
	row := s.db.QueryRow(
		"SELECT id, title, artist, year, genre, notes FROM records WHERE id = ?", id,
	)

	var rec models.Record
	if err := row.Scan(&rec.ID, &rec.Title, &rec.Artist, &rec.Year, &rec.Genre, &rec.Notes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("record not found: %w", err)
		}
		return nil, fmt.Errorf("querying record by ID: %w", err)
	}

	return &rec, nil
}

// Create inserts a new record into the database and returns the created record with its ID.
func (s *RecordService) Create(req models.CreateRecordRequest) (*models.Record, error) {
	result, err := s.db.Exec(
		"INSERT INTO records (title, artist, year, genre, notes) VALUES (?, ?, ?, ?, ?)",
		req.Title, req.Artist, req.Year, req.Genre, req.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("retrieving inserted ID: %w", err)
	}

	return &models.Record{
		ID:     id,
		Title:  req.Title,
		Artist: req.Artist,
		Year:   req.Year,
		Genre:  req.Genre,
		Notes:  req.Notes,
	}, nil
}

// Update applies partial updates to an existing record. It returns the number of
// rows affected so callers can detect whether the record existed.
func (s *RecordService) Update(id int64, req models.UpdateRecordRequest) (int64, error) {
	query := "UPDATE records SET"
	args := []interface{}{}

	fields := 0
	if req.Title != nil {
		fields++
		query += " title = ?"
		args = append(args, *req.Title)
	}
	if req.Artist != nil {
		fields++
		query += ", artist = ?"
		args = append(args, *req.Artist)
	}
	if req.Year != nil {
		fields++
		query += ", year = ?"
		args = append(args, *req.Year)
	}
	if req.Genre != nil {
		fields++
		query += ", genre = ?"
		args = append(args, *req.Genre)
	}
	if req.Notes != nil {
		fields++
		query += ", notes = ?"
		args = append(args, *req.Notes)
	}

	if fields == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("updating record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("retrieving rows affected: %w", err)
	}

	return rowsAffected, nil
}

// Delete removes a record by its ID. It returns the number of rows affected.
func (s *RecordService) Delete(id int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM records WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("deleting record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("retrieving rows affected: %w", err)
	}

	return rowsAffected, nil
}
