package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"records/internal/models"
	"records/internal/services"
)

type Server struct {
	Db *sql.DB
}

func (s *Server) CreateHandler(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /records", func(w http.ResponseWriter, r *http.Request) {

		rs := services.NewRecordService(db)
		recs, err := rs.GetAll()
		if err != nil {
			_ = fmt.Errorf("error when fetching records: %w", err)
		}

		setContentTypeJSON(w)
		_ = json.NewEncoder(w).Encode(recs)
	})

	mux.HandleFunc("POST /records", func(w http.ResponseWriter, r *http.Request) {

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req models.CreateRecordRequest
		err := dec.Decode(&req)
		if err != nil {
			msg := "malformed JSON: " + err.Error()
			writeBadRequestError(w, msg)
			return
		}

	})

	return mux
}

// writeBadRequestError writes a bad request problem details response using errorMsg as the detail.
// It sets the Content-Type to application/problem+json and encodes the response as JSON.
func writeBadRequestError(w http.ResponseWriter, errorMsg string) {
	pd := models.NewProblemDetails(http.StatusBadRequest, errorMsg)
	w.Header().Set("Content-Type", models.ContentTypeProblemJSON)
	_ = json.NewEncoder(w).Encode(&pd)
}

func setContentTypeJSON(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
