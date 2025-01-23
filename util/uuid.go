package util

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UuidString(userId pgtype.UUID) string {
	return uuid.UUID(userId.Bytes).String()
}
