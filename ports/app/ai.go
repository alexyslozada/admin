package app

import (
	"context"

	"github.com/google/uuid"
)

type AI interface {
	CreateThread(ctx context.Context) (uuid.UUID, error)
	CreateMessage(ctx context.Context, threadID uuid.UUID, content string) (string, error)
}
