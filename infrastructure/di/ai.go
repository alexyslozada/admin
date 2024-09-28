package di

import (
	"os"

	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/inbound/http"
	"gitlab.com/EDteam/workshop-ai-2024/admin/adapters/outbound/openai"
	"gitlab.com/EDteam/workshop-ai-2024/admin/application/ai"
)

func InitAIUseCase() ai.UseCase {
	apiKey := os.Getenv("OPENAI_API_KEY")
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")
	if apiKey == "" {
		panic("OPENAI_API_KEY is required")
	}
	if assistantID == "" {
		panic("OPENAI_ASSISTANT_ID is required")
	}

	adapterOpenAI := openai.NewOpenAI(apiKey, assistantID)

	return ai.NewUseCase(adapterOpenAI, InitSaleUseCase())
}

func InitAIHandler() http.AIHandler {
	aiUseCase := InitAIUseCase()
	return http.NewAIHandler(&aiUseCase)
}
