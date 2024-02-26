package ai

import (
	"context"
	"fmt"
	"icelandicicecream/openai-go/model"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client *openai.Client
}

func NewChatClient() *OpenAI {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &OpenAI{
		Client: openai.NewClient(apiKey),
	}
}

func (c *OpenAI) GetCompletion(ctx context.Context, u model.OpenAIRequest, responseChan chan<- string) error {
	var messages []openai.ChatCompletionMessage

	// Add system message
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are an AI assistant. Be courteous and polite when replying.",
	}

	messages = append(messages, systemMessage)

	// Check if there are any previous messages
	if len(u.Messages) != 0 {
		for _, msg := range u.Messages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// Create base request
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 100,
		Messages:  messages,
		// ResponseFormat: &openai.ChatCompletionResponseFormat{
		// 	Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		// },
		Stream: true,
	}

	// Add latest message to the start
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    u.Role,
		Content: u.Content,
	})

	// Get the response stream from openAI
	stream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("Stream error: %v\n", err)
		return err
	}
	defer stream.Close()

	for {
		// Send the response to the channel
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			break
		}

		responseChan <- response.Choices[0].Delta.Content
	}

	return nil
}
