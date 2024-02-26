package server

import (
	"encoding/json"
	"fmt"
	"icelandicicecream/openai-go/model"
	"net/http"
)

func (s *Server) handleRoutes(mux *http.ServeMux, middleware []func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /health", s.applyMiddleware(s.healthCheckHandler, middleware...))
	mux.HandleFunc("POST /openai", s.applyMiddleware(s.sendMessageToOpenAIHandler, middleware...))
}

// applyMiddleware helps in applying a slice of middleware to a handler.
func (s *Server) applyMiddleware(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Handlers
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is up and running! 🚀")
}

func (s *Server) sendMessageToOpenAIHandler(w http.ResponseWriter, r *http.Request) {
	var req model.OpenAIRequest

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a response channel
	responseChan := make(chan string)
	go func() {
		defer close(responseChan)
		err := s.OpenAI.GetCompletion(ctx, req, responseChan)
		if err != nil {
			responseChan <- err.Error()
		}
	}()

	for message := range responseChan {
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher.Flush()
	}

	w.WriteHeader(http.StatusOK)
}
