package bubbletea

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
	"strings"
)

func newGiveOrderToClientModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	const (
		orderIDsInput = iota
	)

	inputs := make([]textinput.Model, 1)

	inputs[orderIDsInput] = textinput.New()
	inputs[orderIDsInput].Focus()
	inputs[orderIDsInput].Prompt = "Order IDs (comma separated): "
	inputs[orderIDsInput].Placeholder = "Enter order ID"

	submit := func(values []string) error {
		orderIDsValue := values[orderIDsInput]

		if orderIDsValue == "" {
			return fmt.Errorf("orderIDs is empty")
		}

		orderIDs := strings.Split(orderIDsValue, ",")
		for i := range orderIDs {
			orderIDs[i] = strings.TrimSpace(orderIDs[i])
		}

		return useCase.GiveOrderToClient(context.Background(), orderIDs)
	}

	return NewFormModel(inputs, submit)
}
