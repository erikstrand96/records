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
	mux.HandleFunc("/", handleRoot)
	return mux
}

func handleRoot(writer http.ResponseWriter, _ *http.Request) {
	setContentTypeHeader(writer)
	_ = json.NewEncoder(writer).Encode(map[string]string{"message": "Hello, this is an application where I keep a collection of all my vinyl records!"})
}

func setContentTypeHeader(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")

}
