package main

import (
	"icelandicicecream/openai-go/pkg/ai"
	"icelandicicecream/openai-go/pkg/db"
	"icelandicicecream/openai-go/pkg/server"
)

func main() {
	// Initialize OpenAI Client
	client := ai.NewChatClient()
	db, err := db.New()
	if err != nil {
		panic(err)
	}

	server := server.Server{
		Config: server.ServerConfig{
			Port: 8080,
		},
		OpenAI: client,
		DB:     db,
	}

	server.Start()
}
