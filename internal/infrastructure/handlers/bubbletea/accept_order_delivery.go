package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"strconv"
	"time"
)

func newAcceptOrderModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	const (
		orderIDInput = iota
		recipientIDInput
		storageTimeInput
		weightInput
		costInput
		packagingInput
		additionalFilmInput
	)

	inputs := make([]textinput.Model, 7)

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

	inputs[weightInput] = textinput.New()
	inputs[weightInput].Prompt = "Weight: "
	inputs[weightInput].Placeholder = "Enter weight"

	inputs[costInput] = textinput.New()
	inputs[costInput].Prompt = "Cost: "
	inputs[costInput].Placeholder = "Enter cost"

	inputs[packagingInput] = textinput.New()
	inputs[packagingInput].Prompt = "Packaging: "
	inputs[packagingInput].Placeholder = "Enter packaging"

	inputs[additionalFilmInput] = textinput.New()
	inputs[additionalFilmInput].Prompt = "Additional film (y/n): "
	inputs[additionalFilmInput].Placeholder = "Enter additional film"

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

		costValue := values[costInput]
		if costValue == "" {
			return fmt.Errorf("cost is empty")
		}

		cost, err := strconv.Atoi(costValue)

		weightValue := values[weightInput]
		if weightValue == "" {
			return fmt.Errorf("weight is empty")
		}

		weight, err := strconv.Atoi(weightValue)

		packagingValue := values[packagingInput]
		if packagingValue == "" {
			return fmt.Errorf("packaging is empty")
		}

		packaging, err := domain.NewPackagingType(packagingValue)
		if err != nil {
			return err
		}

		additionalFilmValue := values[additionalFilmInput]
		if additionalFilmValue == "" {
			return fmt.Errorf("additionalFilm is empty")
		}

		if additionalFilmValue != "y" && additionalFilmValue != "n" {
			return fmt.Errorf("additionalFilm is invalid")
		}

		additionalFilm := additionalFilmValue == "y"

		return useCase.AcceptOrderDelivery(
			orderIDValue, recipientIDValue, storageTime,
			cost, weight, packaging, additionalFilm,
		)
	}

	return NewFormModel(inputs, submit)
}
