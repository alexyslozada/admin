package domain

type (
	AIRunKind      string
	AIFunctionName string
)

const (
	AIRunKindRequiredAction AIRunKind = "required_action"
	AIRunKindResponse       AIRunKind = "response"
)

const (
	AIFunctionNameGetSales AIFunctionName = "GetSales"
)

type AIMessageRequest struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
}

type FunctionCall struct {
	ThreadID string         `json:"thread_id"`
	RunID    string         `json:"run_id"`
	CallID   string         `json:"call_id"`
	Name     string         `json:"name"`
	Args     map[string]any `json:"args"`
}

type Run struct {
	Kind         AIRunKind    `json:"kind"`
	Response     string       `json:"response"`
	FunctionCall FunctionCall `json:"function_call"`
}
