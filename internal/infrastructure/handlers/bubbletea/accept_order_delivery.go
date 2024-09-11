package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
	"time"
)

func newAcceptOrderModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	const (
		orderIDInput = iota
		recipientIDInput
		storageTimeInput
	)

	inputs := make([]textinput.Model, 3)

	inputs[orderIDInput] = textinput.New()
	inputs[orderIDInput].Focus()
	inputs[orderIDInput].Prompt = "Order ID: "
	inputs[orderIDInput].Placeholder = "Enter order ID"

	inputs[recipientIDInput] = textinput.New()
	inputs[recipientIDInput].Prompt = "Recipient ID: "
	inputs[recipientIDInput].Placeholder = "Enter recipient ID"

	inputs[storageTimeInput] = textinput.New()
	inputs[storageTimeInput].Prompt = "Storage time: "
	inputs[storageTimeInput].Placeholder = "Enter storage time"

	submit := func(values []string) error {
		orderIDValue := values[orderIDInput]
		recipientIDValue := values[recipientIDInput]
		storageTimeValue := values[storageTimeInput]

		if orderIDValue == "" {
			return fmt.Errorf("orderID is empty")
		}

		if recipientIDValue == "" {
			return fmt.Errorf("recipientID is empty")
		}

		if storageTimeValue == "" {
			return fmt.Errorf("storageTime is empty")
		}

		storageTime, err := time.ParseDuration(storageTimeValue)
		if err != nil {
			return fmt.Errorf("storageTime is invalid")
		}

		if storageTime < 0 {
			return fmt.Errorf("storageTime is negative")
		}

		return useCase.AcceptOrderDelivery(orderIDValue, recipientIDValue, storageTime)
	}

	return NewFormModel(inputs, submit)
}
