package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"records/internal/config"
)

func main() {

	ctx := context.Background()
	err, cfg := config.NewConfig()
	if err != nil {
		log.Fatalf("Could not create AppConfig: %v", err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		log.Printf(r.Pattern)
		_ = json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Starting server on: %s", addr)
	srv := http.Server{Addr: addr}

	err = srv.ListenAndServe()
	defer func(srv *http.Server, ctx context.Context) {
		_ = srv.Shutdown(ctx)
	}(&srv, ctx)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	}
}
