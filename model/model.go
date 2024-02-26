package model

type OpenAIRequest struct {
	Role     string          `json:"role"`
	Content  string          `json:"content"`
	Messages []OpenAIRequest `json:"messages"`
}
