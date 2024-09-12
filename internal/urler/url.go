package urler

import (
	"net/url"
	"strings"
)

type Operator string

const (
	Equal              Operator = "eq"
	NotEqual           Operator = "ne"
	GreaterThan        Operator = "gt"
	GreaterThanOrEqual Operator = "gte"
	LessThan           Operator = "lt"
	LessThanOrEqual    Operator = "lte"
	Like               Operator = "like"
	In                 Operator = "in"
	NotIn              Operator = "nin"
)

var operators = map[Operator]Void{
	Equal:              {},
	NotEqual:           {},
	GreaterThan:        {},
	GreaterThanOrEqual: {},
	LessThan:           {},
	LessThanOrEqual:    {},
	Like:               {},
	In:                 {},
	NotIn:              {},
}

type Filter struct {
	Field    string
	Operator Operator
	Value    string
}

type Void struct{}

func ParseQueryParams(values url.Values) []Filter {
	var filters []Filter
	for field, valueSlice := range values {
		if len(valueSlice) == 0 {
			continue
		}

		value := valueSlice[0]

		parts := strings.SplitN(value, ":", 2)

		// Default operator is Equal
		// e.g. name=John will return a filter {Field: "name", Operator: Equal, Value: "John"}
		if len(parts) != 2 {
			filter := Filter{
				Field:    field,
				Operator: Equal,
				Value:    value,
			}

			filters = append(filters, filter)
			continue
		}

		operator := parts[0]
		realOperator := Operator(operator)

		// If the operator is not valid, the default operator is Equal
		if _, ok := operators[Operator(operator)]; !ok {
			realOperator = Equal
		}

		actualValue := parts[1]

		filter := Filter{
			Field:    field,
			Operator: realOperator,
			Value:    actualValue,
		}

		filters = append(filters, filter)
	}

	return filters
}
