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
		// Create a buffered channel to keep message ordering
		o.Channels[sessionIdString] = make(chan string, 1)
	}

	return o.Channels[sessionIdString], nil
}

func (o OpenAI) SendCompletion(ctx context.Context, db *db.DB, sessionId pgtype.UUID, message string) error {
	sessionChan, err := o.GetSession(ctx, db, sessionId)
	if err != nil {
		return err
	}

	// Get message history
	messages, err := db.GetSessionMessages(ctx, sessionId)
	if err != nil {
		return err
	}

	// Begin Tx
	tx, rollback, err := db.WithTX(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	// Add User Message
	err = db.AddMessageTx(ctx, tx, sessionId, openai.ChatMessageRoleUser, message)
	if err != nil {
		return err
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	chatReq := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
		Stream:   true,
	}

	chatStream, err := o.Client.CreateChatCompletionStream(ctx, chatReq)
	if err != nil {
		return err
	}
	defer chatStream.Close()

	var chatResponse string

	for {
		if sessionChan == nil {
			fmt.Println("sessionChan is nil")
			return nil
		}
		response, err := chatStream.Recv()
		if errors.Is(err, io.EOF) {

			// Add AI Message to DB
			err = db.AddMessageTx(ctx, tx, sessionId, openai.ChatMessageRoleAssistant, chatResponse)
			if err != nil {
				return err
			}

			if err = tx.Commit(ctx); err != nil {
				return err
			}

			return nil

		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}

		chatResponse += response.Choices[0].Delta.Content
		sessionChan <- response.Choices[0].Delta.Content
	}
}
