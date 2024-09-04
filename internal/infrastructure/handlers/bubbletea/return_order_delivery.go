package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
)

func newReturnOrderModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	const (
		orderIDInput = iota
	)

	inputs := make([]textinput.Model, 1)

	inputs[orderIDInput] = textinput.New()
	inputs[orderIDInput].Focus()
	inputs[orderIDInput].Prompt = "Order ID: "
	inputs[orderIDInput].Placeholder = "Enter order ID"

	submit := func(values []string) error {
		orderIDValue := values[orderIDInput]

		if orderIDValue == "" {
			return fmt.Errorf("orderID is empty")
		}

		return useCase.ReturnOrderDelivery(orderIDValue)
	}

	return NewFormModel(inputs, submit)
}
