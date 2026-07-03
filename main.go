package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"records/internal/config"
	"syscall"
	"time"
)

func main() {

	log.Println("Welcome to records!")

	// Cancel the context when an interrupt/terminate signal arrives so we can
	// shut the server down gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Could not create AppConfig: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		log.Print(r.Pattern)
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := http.Server{Addr: addr, Handler: mux}

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
