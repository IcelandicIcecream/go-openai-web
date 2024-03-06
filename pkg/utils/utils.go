package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgUUID(u uuid.UUID) (pgUUID pgtype.UUID, err error) {
	var byteArray [16]byte
	copy(byteArray[:], u.Bytes())

	pgUUID = pgtype.UUID{
		Valid: true,
		Bytes: byteArray,
	}

	return pgUUID, nil
}

func ConvertUUIDToString(id pgtype.UUID) (string, error) {
	// Convert from pgx UUID type to string
	orgUUID, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		fmt.Println("Error converting UUID:", err)
		return "", err
	}
	return orgUUID.String(), nil
}
