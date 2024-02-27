package server

import (
	"context"
	"fmt"
	"icelandicicecream/openai-go/pkg/ai"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	Port int
}

// Server struct that holds any dependencies for the HTTP server.
type Server struct {
	Config ServerConfig
	OpenAI *ai.OpenAI
}

// Start initializes the HTTP server and its routes.
func (s *Server) Start() {
	mux := http.NewServeMux()

	// Default middleware stack for most routes
	defaultMiddlewares := []func(http.HandlerFunc) http.HandlerFunc{
		s.logRequests,
		s.corsMiddleware,
		s.errorHandler,
	}

	// Handle the routes
	s.handleRoutes(mux, defaultMiddlewares)

	server := &http.Server{
		Addr:         ":" + fmt.Sprint(s.Config.Port),
		Handler:      mux,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the server in a goroutine to handle signals
	go func() {
		fmt.Println("Server is starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", ":8080", err)
		}
	}()

	// Set up signal catching
	sigChan := make(chan os.Signal, 1)
	// Catch all signals since not all signals are captured by default
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// Block until a signal is received
	<-sigChan
	fmt.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	fmt.Println("Server gracefully stopped")
}

// errorHandler is a middleware for error handling.
func (s *Server) errorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from an error: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// corsMiddleware is a middleware to add CORS headers.
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next(w, r)
	}
}

// logRequests is a middleware to log HTTP requests.
func (s *Server) logRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
