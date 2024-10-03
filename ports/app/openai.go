package app

import (
	"context"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
)

type OpenAI interface {
	CreateThread(ctx context.Context) (string, error)
	CreateMessage(ctx context.Context, threadID, content string) (string, error)
	RunThread(ctx context.Context, threadID string) (domain.AIRunKind, []domain.Run, error)
	SubmitToolOutput(ctx context.Context, runners []domain.Run) (domain.AIRunKind, []domain.Run, error)
}
