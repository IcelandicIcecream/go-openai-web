package server

import (
	"fmt"
	"icelandicicecream/openai-go/pkg/ai"
	"icelandicicecream/openai-go/pkg/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type ServerConfig struct {
	Port int
}

// Server struct that holds any dependencies for the HTTP server.
type Server struct {
	Config ServerConfig
	OpenAI *ai.OpenAI
	DB     *db.DB
}

// Start initializes the HTTP server and its routes.
func (s *Server) Start() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("OpenAI Go", "1.0.0"))

	// Default middleware stack for most routes
	defaultMiddlewares := []func(http.Handler) http.Handler{
		s.logRequests,
		s.corsMiddleware,
		// s.errorHandler,
	}

	// Handle the routes
	s.handleRoutes(api)

	server := use(mux, defaultMiddlewares...)

	// Start the server in a goroutine to handle signals
	go func() {
		fmt.Printf("Server is starting on port %v... \n", s.Config.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %d: %v\n", s.Config.Port, err)
		}
	}()

	// Set up signal catching
	sigChan := make(chan os.Signal, 1)
	// Catch all signals since not all signals are captured by default
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// Block until a signal is received
	<-sigChan
	fmt.Println("Shutting down server...")
	s.DB.Close()
	s.OpenAI.Close()

	fmt.Println("Server gracefully stopped")
}

// errorHandler is a middleware for error handling.
// func (s *Server) errorHandler(next http.Handler) http.Handler {
// 	return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Printf("Recovered from an error: %v", err)
// 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			}
// 		}()
// 		next(w, r)
// 	})
// }

// corsMiddleware is a middleware to add CORS headers.
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

// // logRequests is a middleware to log HTTP requests.
func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}
