package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
)

type OpenAI struct {
	apiKey      string
	client      *openai.Client
	assistantID string
}

func NewOpenAI(apiKey string, assistantID string) OpenAI {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
	)
	return OpenAI{apiKey: apiKey, client: client, assistantID: assistantID}
}

func (o OpenAI) CreateThread(ctx context.Context) (string, error) {
	thread, err := o.client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{})
	if err != nil {
		return "", err
	}

	return thread.ID, err
}

func (o OpenAI) CreateMessage(ctx context.Context, threadID, content string) (string, error) {
	message, err := o.client.Beta.Threads.Messages.New(ctx, threadID, openai.BetaThreadMessageNewParams{
		Role: openai.F(openai.BetaThreadMessageNewParamsRoleUser),
		Content: openai.F([]openai.MessageContentPartParamUnion{
			openai.TextContentBlockParam{
				Type: openai.F(openai.TextContentBlockParamTypeText),
				Text: openai.String(content),
			},
		}),
	})
	if err != nil {
		return "", err
	}

	return message.ID, err
}

func (o OpenAI) RunThread(ctx context.Context, threadID string) (domain.AIRunKind, []domain.Run, error) {
	log.Printf("RunThread() threadID: %s", threadID)
	stream := o.client.Beta.Threads.Runs.NewStreaming(
		ctx,
		threadID,
		openai.BetaThreadRunNewParams{
			AssistantID: openai.String(o.assistantID),
		},
	)

	for stream.Next() {
		streamEvent := stream.Current()
		if streamEvent.Event == openai.AssistantStreamEventEventThreadRunCompleted {
			break
		}

		if streamEvent.Event == openai.AssistantStreamEventEventThreadMessageCompleted {
			log.Printf("RunThread() threadID: %s, Message completed", threadID)
			message, ok := streamEvent.Data.(openai.Message)
			if !ok {
				return domain.AIRunKindResponse, nil, fmt.Errorf("could not convert streamEvent.Data to openai.Message, type is: %T", streamEvent.Data)
			}

			responses := make([]domain.Run, 0, len(message.Content))
			for _, content := range message.Content {
				responses = append(responses, domain.Run{Response: content.Text.Value})
			}

			return domain.AIRunKindResponse, responses, nil
		}

		if streamEvent.Event == openai.AssistantStreamEventEventThreadRunRequiresAction {
			log.Printf("RunThread() threadID: %s, Required action", threadID)
			run, ok := streamEvent.Data.(openai.Run)
			if !ok {
				return domain.AIRunKindRequiredAction, nil, fmt.Errorf("could not convert streamEvent.Data to openai.Run, type is: %T", streamEvent.Data)
			}

			runners := make([]domain.Run, 0, len(run.RequiredAction.SubmitToolOutputs.ToolCalls))
			for _, toolCall := range run.RequiredAction.SubmitToolOutputs.ToolCalls {
				log.Printf("RunThread() toolCall: %+v", toolCall)
				function := toolCall.Function

				var args map[string]any
				err := json.Unmarshal([]byte(function.Arguments), &args)
				if err != nil {
					return domain.AIRunKindRequiredAction, nil, err
				}
				log.Printf("RunThread() args: %+v", args)

				functionCall := domain.FunctionCall{
					ThreadID: threadID,
					RunID:    run.ID,
					CallID:   toolCall.ID,
					Name:     function.Name,
					Args:     args,
				}

				runners = append(runners, domain.Run{FunctionCall: functionCall})
			}

			return domain.AIRunKindRequiredAction, runners, nil
		}
	}
	if err := stream.Err(); err != nil {
		return domain.AIRunKindRunCompleted, nil, err
	}

	log.Printf("RunThread() threadID: %s, Run completed, but not data", threadID)
	return domain.AIRunKindRunCompleted, nil, nil
}

func (o OpenAI) SubmitToolOutput(ctx context.Context, runners []domain.Run) (string, error) {
	threadID := runners[0].FunctionCall.ThreadID
	runID := runners[0].FunctionCall.RunID
	outputParams := make([]openai.BetaThreadRunSubmitToolOutputsParamsToolOutput, 0, len(runners))
	for _, runner := range runners {
		outputParams = append(outputParams, openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
			Output:     openai.F(runner.Response),
			ToolCallID: openai.F(runner.FunctionCall.CallID),
		})
	}
	stream := o.client.Beta.Threads.Runs.SubmitToolOutputsStreaming(
		ctx,
		threadID,
		runID,
		openai.BetaThreadRunSubmitToolOutputsParams{ToolOutputs: openai.F(outputParams)},
	)

	for stream.Next() {
		streamEvent := stream.Current()
		if streamEvent.Event == openai.AssistantStreamEventEventThreadRunCompleted {
			break
		}
		if streamEvent.Event == openai.AssistantStreamEventEventThreadMessageCompleted {
			message, ok := streamEvent.Data.(openai.Message)
			if !ok {
				return "", fmt.Errorf("could not convert streamEvent.Data to openai.Message, type is: %T", streamEvent.Data)
			}

			var messageResponse bytes.Buffer
			for _, content := range message.Content {
				messageResponse.WriteString(content.Text.Value)
				messageResponse.WriteString("\n")
			}

			return messageResponse.String(), nil
		}
	}

	return "SubmitToolOutput finish. Data not found", nil
}
