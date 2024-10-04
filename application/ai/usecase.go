package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"gitlab.com/EDteam/workshop-ai-2024/admin/domain"
	"gitlab.com/EDteam/workshop-ai-2024/admin/internal/urler"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/app"
	"gitlab.com/EDteam/workshop-ai-2024/admin/ports/out"
)

type UseCase struct {
	openAI       out.OpenAI
	threads      map[uuid.UUID]string
	salesUseCase app.GenericUseCase[domain.Sale]
}

func NewUseCase(openAI out.OpenAI, sales app.GenericUseCase[domain.Sale]) UseCase {
	return UseCase{
		openAI:       openAI,
		threads:      make(map[uuid.UUID]string),
		salesUseCase: sales,
	}
}

// CreateThread generates a new UUID for a thread, creates a new thread using the OpenAI interface,
// stores the mapping of the generated UUID to the OpenAI thread ID, and returns the generated UUID.
func (uc *UseCase) CreateThread(ctx context.Context) (uuid.UUID, error) {
	// Generate a new UUID for the thread
	ID, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, err
	}

	// Create a new thread using the OpenAI interface
	threadID, err := uc.openAI.CreateThread(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	// Store the mapping of the generated UUID to the OpenAI thread ID
	uc.threads[ID] = threadID

	return ID, nil
}

// CreateMessage retrieves the OpenAI thread ID associated with the provided UUID, creates a new message
// using the OpenAI interface, runs the thread, and returns the response.
func (uc *UseCase) CreateMessage(ctx context.Context, threadID uuid.UUID, content string) (string, error) {
	log.Printf("CreateMessage: threadID=%s, content=%s", threadID, content)
	realThreadID, ok := uc.threads[threadID]
	if !ok {
		return "", domain.ErrThreadNotFound
	}

	_, err := uc.openAI.CreateMessage(ctx, realThreadID, content)
	if err != nil {
		return "", err
	}

	runID, err := uc.openAI.RunThread(ctx, realThreadID)
	if err != nil {
		return "", err
	}

	err = uc.getRunCompleted(ctx, realThreadID, runID)
	if err != nil {
		return "", err
	}

	messages, err := uc.openAI.GetMessagesFromRun(ctx, realThreadID, runID)
	if err != nil {
		return "", err
	}

	response := strings.Join(messages, "\n")
	return response, nil
}

func (uc *UseCase) getRunCompleted(ctx context.Context, threadID, runID string) error {
	var status domain.AIRunKind
	var requiredAction domain.AIRequiredAction
	var runners []domain.Run
	var err error

	// Poll the OpenAI interface for the run status
	// try 50 times with a 2-second delay between each attempt
	for range 50 {
		time.Sleep(2 * time.Second)

		status, requiredAction, runners, err = uc.openAI.GetRun(ctx, threadID, runID)
		if err != nil {
			return err
		}

		if status == domain.AIRunKindRunCompleted {
			return nil
		}

		if status == domain.AIRunKindRequiresAction && requiredAction == domain.AIRequiredActionSubmitToolOutputs {
			for i, runner := range runners {
				if runner.FunctionCall.Name == domain.AIFunctionNameGetSales {
					// Perform the GetSales action
					actionResponse, err := uc.performAction(ctx, runner)
					if err != nil {
						log.Printf("GetRunCompleted() performAction() error: %v", err)
						return err
					}

					runners[i].Response = actionResponse
				}
			}

			err = uc.openAI.SubmitToolOutput(ctx, threadID, runID, runners)
			if err != nil {
				return err
			}
		}

		log.Printf("GetRunCompleted() status: %s, requiredAction: %s", status, requiredAction)
	}

	log.Printf("GetRunCompleted() run did not complete")
	return errors.New("run did not complete")
}

func (uc *UseCase) performAction(ctx context.Context, run domain.Run) (string, error) {
	// Perform the required action
	if domain.AIFunctionName(run.FunctionCall.Name) == domain.AIFunctionNameGetSales {
		from, ok := run.FunctionCall.Args["from"].(string)
		if !ok {
			return "", fmt.Errorf("could not convert 'from' argument to string")
		}
		to, ok := run.FunctionCall.Args["to"].(string)
		if !ok {
			return "", fmt.Errorf("could not convert 'to' argument to string")
		}
		kind, _ := run.FunctionCall.Args["kind"].(string)
		filter := []urler.Filter{
			{
				Field:    "from",
				Operator: urler.Equal,
				Value:    from,
			},
			{
				Field:    "to",
				Operator: urler.Equal,
				Value:    to,
			},
		}
		if kind != "" {
			filter = append(filter, urler.Filter{
				Field:    "kind",
				Operator: urler.Equal,
				Value:    kind,
			})
		}

		// Perform the GetSales action
		sales, err := uc.salesUseCase.FindAll(filter)
		if err != nil {
			return "", err
		}

		salesDTO := make([]domain.SaleResponse, 0, len(sales))
		for _, sale := range sales {
			salesDTO = append(salesDTO, domain.SaleResponse{
				ID:             sale.ID,
				Product:        sale.Product,
				ClientID:       sale.ClientID,
				Client:         sale.Client,
				DateInvoice:    time.Unix(sale.DateInvoice, 0),
				Amount:         sale.Amount,
				IsSubscription: sale.IsSubscription,
				Months:         sale.Months,
				CreatedAt:      time.Unix(sale.CreatedAt, 0),
				UpdatedAt:      time.Unix(sale.UpdatedAt, 0),
				DeletedAt:      time.Unix(sale.DeletedAt, 0),
			})
		}

		raw, err := json.Marshal(salesDTO)
		if err != nil {
			return "", err
		}

		return string(raw), nil
	}

	return "", fmt.Errorf("unknown function name: %s", run.FunctionCall.Name)
}
