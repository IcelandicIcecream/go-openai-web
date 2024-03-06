package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type DB struct {
	Pool *pgxpool.Pool
}

var (
	addSessionSQL         = "INSERT INTO sessions (user_id) VALUES ($1) RETURNING id"
	checkSessionExistSQL  = "SELECT EXISTS (SELECT id FROM sessions WHERE id = $1)"
	addMessageSQL         = "INSERT INTO messages (session_id, content, role) VALUES ($1, $2, $3)"
	getSessionMessagesSQL = "SELECT content, role FROM messages WHERE session_id = $1"
)

func New() (*DB, error) {
	// load env
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	pgUrl := os.Getenv("PG_URL")
	poolConfig, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &DB{
		Pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Ping() error {
	return db.Pool.Ping(context.Background())
}

func (db *DB) WithTX(ctx context.Context) (pgx.Tx, func(), error) {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	return tx, func() { _ = tx.Rollback(ctx) }, nil // Rollback if not committed
}

func (db *DB) AddSession(ctx context.Context, userId pgtype.UUID) (pgtype.UUID, error) {
	var sessionId pgtype.UUID
	err := db.Pool.QueryRow(ctx, addSessionSQL, userId).Scan(&sessionId)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return sessionId, nil
}

func (db *DB) CheckSessionExists(ctx context.Context, sessionId pgtype.UUID) (pgtype.Bool, error) {
	var sessionExists pgtype.Bool

	err := db.Pool.QueryRow(ctx, checkSessionExistSQL, sessionId).Scan(&sessionExists)
	if err != nil {
		return pgtype.Bool{}, err
	}

	return sessionExists, nil
}

func (db *DB) GetSessionMessages(ctx context.Context, sessionId pgtype.UUID) (messages []openai.ChatCompletionMessage, err error) {
	rows, err := db.Pool.Query(ctx, getSessionMessagesSQL, sessionId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message openai.ChatCompletionMessage
		if err := rows.Scan(&message.Content, &message.Role); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (db *DB) AddMessageTx(ctx context.Context, tx pgx.Tx, sessionId pgtype.UUID, role, message string) error {
	_, err := tx.Exec(ctx, addMessageSQL, sessionId, message, role)
	if err != nil {
		return err
	}
	return err
}
