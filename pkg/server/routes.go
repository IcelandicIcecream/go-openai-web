package server

import (
	"context"
	"fmt"
	"icelandicicecream/openai-go/model"
	"icelandicicecream/openai-go/pkg/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

func (s *Server) handleRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		OperationID:   "health-check",
		Path:          "/health",
		Summary:       "Check server health",
		DefaultStatus: http.StatusOK,
	}, s.healthCheckHandler)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		OperationID:   "new-chat-session",
		Path:          "/openai/chat",
		Summary:       "Start conversation with OpenAI",
		DefaultStatus: http.StatusOK,
	}, s.newSession)

	sse.Register(api, huma.Operation{
		Method:        http.MethodGet,
		OperationID:   "connect-chat-session",
		Path:          "/openai/chat/{session_id}",
		Summary:       "Stream conversation from OpenAI",
		DefaultStatus: http.StatusOK,
	}, map[string]any{
		"stream": model.StreamResponse{},
	}, s.streamSession)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		OperationID:   "send-chat-completion",
		Path:          "/openai/chat/complete",
		Summary:       "Send text completions to OpenAI",
		DefaultStatus: http.StatusOK,
	}, s.chatCompletion)
}

// Register wraps `huma.Register` to provide automatic error conversion.
func Register[I, O any](api huma.API, op huma.Operation, handler func(context.Context, *I) (*O, error)) {
	huma.Register(api, op, func(ctx context.Context, input *I) (*O, error) {
		return handler(ctx, input)
	})
}

func (s *Server) healthCheckHandler(ctx context.Context, input *struct{}) (*model.Response, error) {
	resp := &model.Response{}
	resp.Body.Message = "Server is up and running! 🚀"
	return resp, nil
}

// create a new session
func (s *Server) newSession(ctx context.Context, req *model.NewSessionRequest) (*model.Response, error) {
	userId, err := utils.PgUUID(req.Body.UserId)
	if err != nil {
		return nil, err
	}

	sessionId, err := s.OpenAI.NewSession(ctx, s.DB, userId)
	if err != nil {
		return nil, err
	}

	resp := &model.Response{}
	resp.Body.Message = "Session created successfully"
	resp.Body.Payload = sessionId

	return resp, nil
}

// connect to session
func (s *Server) streamSession(ctx context.Context, req *model.StreamSessionRequest, send sse.Sender) {
	sessionId, err := utils.PgUUID(req.SessionId)
	if err != nil {
		return
	}
	sessionChan, err := s.OpenAI.GetSession(ctx, s.DB, sessionId)
	if err != nil {
		return
	}
	for message := range sessionChan {
		fmt.Println(message)
		send.Data(model.StreamResponse{Message: message})
	}
}

func (s *Server) chatCompletion(ctx context.Context, req *model.OpenAIRequest) (*struct{}, error) {
	sessionId, err := utils.PgUUID(req.Body.SessionId)
	if err != nil {
		return nil, err
	}

	err = s.OpenAI.SendCompletion(ctx, s.DB, sessionId, req.Body.Message)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
