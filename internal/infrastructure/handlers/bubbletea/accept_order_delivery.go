package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"strconv"
	"time"
)

const (
	acceptOrderModelOrderIDInput = iota
	acceptOrderModelRecipientIDInput
	acceptOrderModelStorageTimeInput
	acceptOrderModelWeightInput
	acceptOrderModelCostInput
	acceptOrderModelPackagingInput
	acceptOrderModelAdditionalFilmInput
)

func initInputs() []textinput.Model {
	inputs := make([]textinput.Model, 7)

	inputs[acceptOrderModelOrderIDInput] = textinput.New()
	inputs[acceptOrderModelOrderIDInput].Focus()
	inputs[acceptOrderModelOrderIDInput].Prompt = "Order ID: "
	inputs[acceptOrderModelOrderIDInput].Placeholder = "Enter order ID"

	inputs[acceptOrderModelRecipientIDInput] = textinput.New()
	inputs[acceptOrderModelRecipientIDInput].Prompt = "Recipient ID: "
	inputs[acceptOrderModelRecipientIDInput].Placeholder = "Enter recipient ID"

	inputs[acceptOrderModelStorageTimeInput] = textinput.New()
	inputs[acceptOrderModelStorageTimeInput].Prompt = "Storage time: "
	inputs[acceptOrderModelStorageTimeInput].Placeholder = "Enter storage time"

	inputs[acceptOrderModelWeightInput] = textinput.New()
	inputs[acceptOrderModelWeightInput].Prompt = "Weight: "
	inputs[acceptOrderModelWeightInput].Placeholder = "Enter weight"

	inputs[acceptOrderModelCostInput] = textinput.New()
	inputs[acceptOrderModelCostInput].Prompt = "Cost: "
	inputs[acceptOrderModelCostInput].Placeholder = "Enter cost"

	inputs[acceptOrderModelPackagingInput] = textinput.New()
	inputs[acceptOrderModelPackagingInput].Prompt = "Packaging: "
	inputs[acceptOrderModelPackagingInput].Placeholder = "Enter packaging"

	inputs[acceptOrderModelAdditionalFilmInput] = textinput.New()
	inputs[acceptOrderModelAdditionalFilmInput].Prompt = "Additional film (y/n): "
	inputs[acceptOrderModelAdditionalFilmInput].Placeholder = "Enter additional film"

	return inputs
}

type inputValues struct {
	OrderID        string
	RecipientID    string
	StorageTime    string
	Weight         string
	Cost           string
	Packaging      string
	AdditionalFilm string
}

type validatedInputValues struct {
	OrderID        string
	RecipientID    string
	StorageTime    time.Duration
	Weight         int
	Cost           int
	Packaging      domain.PackagingType
	AdditionalFilm bool
}

func validateOrderID(orderID string) (string, error) {
	if orderID == "" {
		return "", fmt.Errorf("orderID is empty")
	}

	return orderID, nil
}

func validateRecipientID(recipientID string) (string, error) {
	if recipientID == "" {
		return "", fmt.Errorf("recipientID is empty")
	}

	return recipientID, nil
}

func validateStorageTime(storageTime string) (time.Duration, error) {
	if storageTime == "" {
		return 0, fmt.Errorf("storageTime is empty")
	}

	return time.ParseDuration(storageTime)
}

func validateWeight(weight string) (int, error) {
	if weight == "" {
		return 0, fmt.Errorf("weight is empty")
	}

	return strconv.Atoi(weight)
}

func validateCost(cost string) (int, error) {
	if cost == "" {
		return 0, fmt.Errorf("cost is empty")
	}

	return strconv.Atoi(cost)
}

func validatePackaging(packaging string) (domain.PackagingType, error) {
	if packaging == "" {
		return domain.PackagingTypeUnknown, fmt.Errorf("packaging is empty")
	}

	return domain.NewPackagingType(packaging)
}

func validateAdditionalFilm(additionalFilm string) (bool, error) {
	if additionalFilm == "" {
		return false, fmt.Errorf("additionalFilm is empty")
	}

	if additionalFilm != "y" && additionalFilm != "n" {
		return false, fmt.Errorf("additionalFilm is invalid")
	}

	return additionalFilm == "y", nil
}

func validateInputValues(input inputValues) (validatedInputValues, error) {
	errs := make([]error, 0, 7)

	var err error
	var validated validatedInputValues

	validated.OrderID, err = validateOrderID(input.OrderID)
	errs = append(errs, err)

	validated.RecipientID, err = validateRecipientID(input.RecipientID)
	errs = append(errs, err)

	validated.StorageTime, err = validateStorageTime(input.StorageTime)
	errs = append(errs, err)

	validated.Weight, err = validateWeight(input.Weight)
	errs = append(errs, err)

	validated.Cost, err = validateCost(input.Cost)
	errs = append(errs, err)

	validated.Packaging, err = validatePackaging(input.Packaging)
	errs = append(errs, err)

	validated.AdditionalFilm, err = validateAdditionalFilm(input.AdditionalFilm)
	errs = append(errs, err)

	for _, err := range errs {
		if err != nil {
			return validated, err
		}
	}

	return validated, nil
}

func acceptOrderModelSubmit(useCase abstractions.IPVZOrderUseCase) func(values []string) error {
	return func(values []string) error {
		input := inputValues{
			OrderID:        values[acceptOrderModelOrderIDInput],
			RecipientID:    values[acceptOrderModelRecipientIDInput],
			StorageTime:    values[acceptOrderModelStorageTimeInput],
			Weight:         values[acceptOrderModelWeightInput],
			Cost:           values[acceptOrderModelCostInput],
			Packaging:      values[acceptOrderModelPackagingInput],
			AdditionalFilm: values[acceptOrderModelAdditionalFilmInput],
		}

		validated, err := validateInputValues(input)
		if err != nil {
			return err
		}

		return useCase.AcceptOrderDelivery(
			validated.OrderID, validated.RecipientID, validated.StorageTime,
			validated.Cost, validated.Weight, validated.Packaging, validated.AdditionalFilm,
		)
	}
}

func newAcceptOrderModel(useCase abstractions.IPVZOrderUseCase) *FormModel {
	inputs := initInputs()

	submit := acceptOrderModelSubmit(useCase)

	return NewFormModel(inputs, submit)
}
