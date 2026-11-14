package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"records/internal/models"
)

type Server struct {
	Db *sql.DB
}

func (s *Server) CreateRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /records", func(w http.ResponseWriter, r *http.Request) {
		getAllRecords(w, r, db)
	})
	return mux
}

func getAllRecords(writer http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("select id, title, artist, year, genre, notes from records")
	if err != nil {
		log.Printf("error querying records: %v\n", err)
		prob := models.NewProblemDetails(http.StatusInternalServerError, "Failed to query records from database").
			WithInstance(r.URL.Path)
		_ = prob.Write(writer)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("error closing database connection: %v\n", err)
		}
	}(rows)

	var records []models.Record
	for rows.Next() {
		var rec models.Record
		if err := rows.Scan(&rec.ID, &rec.Title, &rec.Artist, &rec.Year, &rec.Genre, &rec.Notes); err != nil {
			log.Printf("error scanning record: %v\n", err)
			prob := models.NewProblemDetails(http.StatusInternalServerError, "Failed to read record").
				WithInstance(r.URL.Path)
			_ = prob.Write(writer)
			return
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error iterating records: %v\n", err)
		prob := models.NewProblemDetails(http.StatusInternalServerError, "Error iterating records").
			WithInstance(r.URL.Path)
		_ = prob.Write(writer)
		return
	}

	if records == nil {
		records = []models.Record{}
	}

	setContentTypeJSON(writer)
	_ = json.NewEncoder(writer).Encode(records)
}

func setContentTypeJSON(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
