package out

import (
	"context"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
)

type OpenAI interface {
	CreateThread(ctx context.Context) (string, error)
	CreateMessage(ctx context.Context, threadID, content string) (string, error)
	RunThread(ctx context.Context, threadID string) (string, error)
	GetRun(ctx context.Context, threadID, runID string) (domain.AIRunKind, domain.AIRequiredAction, []domain.Run, error)
	GetMessagesFromRun(ctx context.Context, threadID, runID string) ([]string, error)
	SubmitToolOutput(ctx context.Context, threadID, runID string, runners []domain.Run) error
}
