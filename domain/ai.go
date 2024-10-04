package domain

type (
	AIRunKind      string
	AIFunctionName string
)

const (
	AIRunKindRequiresAction AIRunKind = "requires_action"
	AIRunKindResponse       AIRunKind = "response"
	AIRunKindRunCompleted   AIRunKind = "completed"
)

const (
	AIFunctionNameGetSales AIFunctionName = "GetSales"
)

type AIMessageRequest struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
}

type FunctionCall struct {
	CallID   string         `json:"call_id"`
	Name     AIFunctionName `json:"name"`
	Args     map[string]any `json:"args"`
	Response string         `json:"response"`
}

type Run struct {
	Response     string       `json:"response"`
	FunctionCall FunctionCall `json:"function_call"`
}
