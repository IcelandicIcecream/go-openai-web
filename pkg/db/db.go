package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DB struct {
	Pool *pgxpool.Pool
}

var (
	addSessionSQL        = "INSERT INTO sessions (user_id) VALUES ($1) RETURNING id"
	checkSessionExistSQL = "SELECT EXISTS (SELECT id FROM sessions WHERE id = $1)"
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
