package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"records/internal/config"
	"records/internal/database"
	"records/internal/models"
	"syscall"
	"time"
)

func main() {

	// Cancel the context when an interrupt/terminate signal arrives so we can
	// shut the server down gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("RECORDS_PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("RECORDS_DB_PATH")
	if dbPath == "" {
		dbPath = "records.db"
	}

	cfg, err := config.NewConfig(host, port, dbPath)
	if err != nil {
		log.Fatalf("Could not create AppConfig: %v", err)
	}

	// Open the SQLite database and create the schema.
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	router := createRouter(db)
	srv := http.Server{Addr: addr, Handler: router}

	// Run the server in a goroutine so main can wait on the shutdown signal.
	go func() {
		log.Printf("Starting server on: %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Block until a signal is received, then shut down gracefully.
	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

func createRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleRoot)

	mux.HandleFunc("GET /records", func(w http.ResponseWriter, r *http.Request) {
		getAllRecords(w, r, db)
	})
	return mux
}

func handleRoot(writer http.ResponseWriter, _ *http.Request) {
	setContentTypeJSON(writer)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"message": "Hello, this is an application where I keep a collection of all my vinyl records!"})
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

func setContentTypeProblemJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", models.ContentTypeProblemJSON)
}
