package model

import (
	"github.com/gofrs/uuid"
)

type OpenAICompletion struct {
	Message   string    `json:"message"`
	SessionId uuid.UUID `json:"session_id" path:"session_id"`
}

type OpenAIRequest struct {
	Body OpenAICompletion `json:"body"`
}

type AccountingResponse struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

type StreamResponse struct {
	Message string `json:"message"`
}

type StreamSession struct {
	SessionId uuid.UUID `json:"session_id" path:"session_id"`
}

type StreamSessionRequest struct {
	SessionId uuid.UUID `json:"session_id" path:"session_id"`
}

type Session struct {
	UserId    string    `json:"user_id"`
	SessionId uuid.UUID `json:"session_id,omitempty"`
}

type NewSessionRequest struct {
	Body Session `json:"body"`
}

type Response struct {
	Body JSONResponse `json:"body"`
}

type JSONResponse struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}
