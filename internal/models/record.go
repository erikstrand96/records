package models

// Record represents a vinyl record in the collection.
type Record struct {
	ID     int64   `json:"id" db:"id"`
	Title  string  `json:"title" db:"title"`
	Artist string  `json:"artist" db:"artist"`
	Year   *int    `json:"year,omitempty" db:"year"`
	Genre  *string `json:"genre,omitempty" db:"genre"`
	Notes  *string `json:"notes,omitempty" db:"notes"`
}

// CreateRecordRequest represents the payload used to create a new record.
type CreateRecordRequest struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Year   *int    `json:"year,omitempty"`
	Genre  *string `json:"genre,omitempty"`
	Notes  *string `json:"notes,omitempty"`
}

// UpdateRecordRequest represents the payload used to update an existing record.
type UpdateRecordRequest struct {
	Title  *string `json:"title,omitempty"`
	Artist *string `json:"artist,omitempty"`
	Year   *int    `json:"year,omitempty"`
	Genre  *string `json:"genre,omitempty"`
	Notes  *string `json:"notes,omitempty"`
}
