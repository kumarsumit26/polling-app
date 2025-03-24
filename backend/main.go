package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"go-api-project/internal/api"
	"go-api-project/internal/repository"
)

func main() {
	repo, err := repository.ConnectPostgresDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	defer repo.Close()

	// Set up the API routes
	router := api.NewRouter(repo)

	// Create a CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                             // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Allow specific methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allow specific headers
		AllowCredentials: true,
	})

	// Wrap your router with the CORS handler
	handler := c.Handler(loggingMiddleware(router))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// loggingMiddleware logs the success or failure of each API call
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		log.Printf("Method: %s, URL: %s, Status: %d, Duration: %s", r.Method, r.URL.Path, ww.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
