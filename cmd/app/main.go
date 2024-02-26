package main

import (
	"icelandicicecream/openai-go/pkg/ai"
	"icelandicicecream/openai-go/pkg/server"
)

func main() {
	// Initialize OpenAI Client
	client := ai.NewChatClient()

	server := server.Server{
		Config: server.ServerConfig{
			Port: 8080,
		},
		OpenAI: client,
	}

	server.Start()
}
