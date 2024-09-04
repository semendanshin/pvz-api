package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
)

func newAcceptReturnModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	const (
		recipientIDInput = iota
		orderIDInput
	)

	inputs := make([]textinput.Model, 2)

	inputs[recipientIDInput] = textinput.New()
	inputs[recipientIDInput].Focus()
	inputs[recipientIDInput].Prompt = "Recipient ID: "
	inputs[recipientIDInput].Placeholder = "Enter recipient ID"

	inputs[orderIDInput] = textinput.New()
	inputs[orderIDInput].Prompt = "Order ID: "
	inputs[orderIDInput].Placeholder = "Enter order ID"

	submit := func(values []string) error {
		recipientIDValue := values[recipientIDInput]
		orderIDValue := values[orderIDInput]

		if recipientIDValue == "" {
			return fmt.Errorf("recipientID is empty")
		}

		if orderIDValue == "" {
			return fmt.Errorf("orderID is empty")
		}

		return useCase.AcceptReturn(recipientIDValue, orderIDValue)
	}

	return NewFormModel(inputs, submit)
}
