package ai

import (
	"context"
	"errors"
	"fmt"
	"icelandicicecream/openai-go/pkg/db"
	"icelandicicecream/openai-go/pkg/utils"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client   *openai.Client
	Channels map[string]chan string
}

var orgSchema = "90d3c048-f545-4972-871c-64c383eccb0d"

func NewChatClient() *OpenAI {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	channels := make(map[string]chan string)

	return &OpenAI{
		Client:   openai.NewClient(apiKey),
		Channels: channels,
	}
}

func (o OpenAI) Close() {
	for _, ch := range o.Channels {
		close(ch)
	}
}

func (o OpenAI) NewSession(ctx context.Context, db *db.DB, userId pgtype.UUID) (sessionId pgtype.UUID, err error) {
	// Insert into DB and get sessionId
	sessionId, err = db.AddSession(ctx, userId)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return sessionId, nil
}

func (o OpenAI) GetSession(ctx context.Context, db *db.DB, sessionId pgtype.UUID) (chan string, error) {
	fmt.Println("SESSION ID: ", sessionId)
	// Check if session exists
	exists, err := db.CheckSessionExists(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	if !exists.Bool {
		return nil, err
	}

	sessionIdString, err := utils.ConvertUUIDToString(sessionId)
	if err != nil {
		return nil, err
	}

	// Create the channel if its not already in memory
	if _, ok := o.Channels[sessionIdString]; !ok {
		o.Channels[sessionIdString] = make(chan string)
	}

	return o.Channels[sessionIdString], nil
}

func (o OpenAI) SendCompletion(ctx context.Context, db *db.DB, sessionId pgtype.UUID, message string) error {
	sessionChan, err := o.GetSession(ctx, db, sessionId)
	if err != nil {
		return err
	}

	chatReq := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		Stream: true,
	}

	chatStream, err := o.Client.CreateChatCompletionStream(ctx, chatReq)
	if err != nil {
		return err
	}
	defer chatStream.Close()

	for {
		if sessionChan == nil {
			fmt.Println("sessionChan is nil")
			return nil
		}
		response, err := chatStream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return nil
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}
		sessionChan <- response.Choices[0].Delta.Content
	}
}
