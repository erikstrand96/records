package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"records/internal/config"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg, err := config.NewConfig(host, port)
	if err != nil {
		log.Fatalf("Could not create AppConfig: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	router := createRouter()
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

func createRouter() *http.ServeMux {
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
