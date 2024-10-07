package openai

import (
	"context"
	"encoding/json"
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

func (o OpenAI) RunThread(ctx context.Context, threadID string) (string, error) {
	log.Printf("RunThread() threadID: %s", threadID)
	run, err := o.client.Beta.Threads.Runs.New(
		ctx,
		threadID,
		openai.BetaThreadRunNewParams{
			AssistantID: openai.String(o.assistantID),
		},
	)

	if err != nil {
		log.Printf("RunThread() error: %s", err)
		return "", err
	}

	return run.ID, nil
}

func (o OpenAI) GetRun(ctx context.Context, threadID, runID string) (domain.AIRunKind, domain.AIRequiredAction, []domain.Run, error) {
	log.Printf("GetRun() threadID: %s, runID: %s", threadID, runID)
	run, err := o.client.Beta.Threads.Runs.Get(ctx, threadID, runID)
	if err != nil {
		log.Printf("GetRun() error: %v", err)
		return "", "", nil, err
	}

	var runners []domain.Run
	for _, toolCall := range run.RequiredAction.SubmitToolOutputs.ToolCalls {
		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			log.Printf("GetRun() error: %v", err)
			return "", "", nil, err
		}
		log.Printf("GetRun() function name: %s, args: %+v", toolCall.Function.Name, args)
		runner := domain.Run{
			FunctionCall: domain.FunctionCall{
				CallID: toolCall.ID,
				Name:   domain.AIFunctionName(toolCall.Function.Name),
				Args:   args,
			},
		}
		runners = append(runners, runner)
	}

	return domain.AIRunKind(run.Status), domain.AIRequiredAction(run.RequiredAction.Type), runners, nil
}

func (o OpenAI) GetMessagesFromRun(ctx context.Context, threadID, runID string) ([]string, error) {
	log.Printf("GetMessagesFromRun() threadID: %s, runID: %s", threadID, runID)
	message, err := o.client.Beta.Threads.Messages.List(ctx, threadID, openai.BetaThreadMessageListParams{RunID: openai.F(runID)})
	if err != nil {
		log.Printf("GetMessagesFromRun() error: %v", err)
		return nil, err
	}

	responses := make([]string, 0, len(message.Data))
	for _, message := range message.Data {
		for _, content := range message.Content {
			responses = append(responses, content.Text.Value)
		}
	}

	return responses, nil
}

func (o OpenAI) SubmitToolOutput(ctx context.Context, threadID, runID string, runners []domain.Run) error {
	outputParams := make([]openai.BetaThreadRunSubmitToolOutputsParamsToolOutput, 0, len(runners))
	for _, runner := range runners {
		outputParams = append(outputParams, openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
			Output:     openai.F(runner.Response),
			ToolCallID: openai.F(runner.FunctionCall.CallID),
		})
	}
	log.Printf("SubmitToolOutput() outputParams: %+v", outputParams)

	run, err := o.client.Beta.Threads.Runs.SubmitToolOutputs(
		ctx,
		threadID,
		runID,
		openai.BetaThreadRunSubmitToolOutputsParams{ToolOutputs: openai.F(outputParams)},
	)
	if err != nil {
		log.Printf("SubmitToolOutput() error: %v", err)
		return err
	}

	log.Printf("SubmitToolOutput() threadID: %s, runID: %s, status: %s", threadID, runID, run.Status)
	return nil
}
