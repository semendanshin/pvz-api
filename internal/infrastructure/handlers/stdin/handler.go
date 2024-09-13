package stdin

import (
	"bufio"
	"fmt"
	"homework/internal/abstractions"
	"os"
	"strconv"
	"strings"
	"time"
)

// Command is a type for command
type Command string

// CommandInput is a type for command input
type CommandInput interface{}

// Command constants
const (
	AcceptOrderCommand  = "accept_delivery"
	ReturnOrderCommand  = "return_delivery"
	GiveOrderCommand    = "give_order"
	GetOrdersCommand    = "get_orders"
	AcceptReturnCommand = "accept_return"
	GetReturnsCommand   = "get_returns"
)

// AcceptOrderCommandInput is a struct for accept order command input
type AcceptOrderCommandInput struct {
	OrderID     string
	RecipientID string
	StorageTime time.Duration
}

// ReturnOrderCommandInput is a struct for return order command input
type ReturnOrderCommandInput struct {
	OrderID string
}

// GiveOrderCommandInput is a struct for give order command input
type GiveOrderCommandInput struct {
	OrderIDs []string
}

// GetOrdersCommandInput is a struct for get orders command input
type GetOrdersCommandInput struct {
	UserID string
}

// AcceptReturnCommandInput is a struct for accept return command input
type AcceptReturnCommandInput struct {
	UserID  string
	OrderID string
}

// GetReturnsCommandInput is a struct for get returns command input
type GetReturnsCommandInput struct {
}

// Handler is a struct for handler
type Handler struct {
	useCase abstractions.IPVZOrderUseCase
}

// NewHandler is a constructor for Handler
func NewHandler(useCase abstractions.IPVZOrderUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Run runs the handler
func (h *Handler) Run() error {
	var command Command
	var input CommandInput
	var inputLine, commandText string

	reader := bufio.NewReader(os.Stdin)

	_, err := fmt.Println("Enter command:")
	if err != nil {
		return err
	}

	for {
		inputLine, _ = reader.ReadString('\n')
		args := strings.Split(strings.Trim(inputLine, "\n"), " ")

		commandText, args = args[0], args[1:]

		switch commandText {
		case AcceptOrderCommand:
			if len(args) < 3 {
				_, err := fmt.Println("Not enough arguments\nUsage: accept_delivery <order_id> <recipient_id> <storage_time>")
				if err != nil {
					return err
				}
				continue
			}

			storageTime, err := strconv.Atoi(args[2])
			if err != nil {
				_, err := fmt.Println("Invalid storage time")
				if err != nil {
					return err
				}
				continue
			}

			if storageTime < 0 {
				_, err := fmt.Println("Storage time is negative")
				if err != nil {
					return err
				}
				continue
			}

			command = AcceptOrderCommand
			input = AcceptOrderCommandInput{
				OrderID:     args[0],
				RecipientID: args[1],
				StorageTime: time.Duration(storageTime) * time.Hour,
			}
		case ReturnOrderCommand:
			if len(args) < 1 {
				_, err := fmt.Println("Not enough arguments\nUsage: return_delivery <order_id>")
				if err != nil {
					return err
				}
				continue
			}

			command = ReturnOrderCommand
			input = ReturnOrderCommandInput{
				OrderID: args[0],
			}
		case GiveOrderCommand:
			if len(args) < 1 {
				_, err := fmt.Println("Not enough arguments\nUsage: give_order <order_id1> <order_id2> ...")
				if err != nil {
					return err
				}
				continue
			}

			command = GiveOrderCommand
			input = GiveOrderCommandInput{
				OrderIDs: args,
			}
		case GetOrdersCommand:
			if len(args) < 1 {
				_, err := fmt.Println("Not enough arguments\nUsage: get_orders <user_id>")
				if err != nil {
					return err
				}
				continue
			}
			command = GetOrdersCommand
			input = GetOrdersCommandInput{
				UserID: args[0],
			}
		case AcceptReturnCommand:
			if len(args) < 2 {
				_, err := fmt.Println("Not enough arguments\nUsage: accept_return <user_id> <order_id>")
				if err != nil {
					return err
				}
				continue
			}

			command = AcceptReturnCommand
			input = AcceptReturnCommandInput{
				UserID:  args[0],
				OrderID: args[1],
			}
		case GetReturnsCommand:
			command = GetReturnsCommand
			input = GetReturnsCommandInput{}
		default:
			_, err := fmt.Println("Unknown command")
			if err != nil {
				return err
			}
		}

		err = h.Handle(command, input)
		if err != nil {
			_, err := fmt.Println("Error accepting order:", err)
			if err != nil {
				return err
			}
		}
	}
}

// Handle handles the command
func (h *Handler) Handle(command Command, input CommandInput) error {
	switch command {
	case AcceptOrderCommand:
		input := input.(AcceptOrderCommandInput)
		err := h.useCase.AcceptOrderDelivery(input.OrderID, input.RecipientID, input.StorageTime)
		if err != nil {
			return err
		}
		_, err = fmt.Println("Order accepted")
		if err != nil {
			return err
		}
	case ReturnOrderCommand:
		input := input.(ReturnOrderCommandInput)
		err := h.useCase.ReturnOrderDelivery(input.OrderID)
		if err != nil {
			return err
		}
		_, err = fmt.Println("Order returned")
		if err != nil {
			return err
		}
	case GiveOrderCommand:
		input := input.(GiveOrderCommandInput)
		err := h.useCase.GiveOrderToClient(input.OrderIDs)
		if err != nil {
			return err
		}
		_, err = fmt.Println("Order given")
		if err != nil {
			return err
		}
	case GetOrdersCommand:
		input := input.(GetOrdersCommandInput)
		data, err := h.useCase.GetOrders(input.UserID)
		if err != nil {
			return err
		}
		for _, order := range data {
			_, err := fmt.Println(order)
			if err != nil {
				return err
			}
		}
	case AcceptReturnCommand:
		input := input.(AcceptReturnCommandInput)
		err := h.useCase.AcceptReturn(input.UserID, input.OrderID)
		if err != nil {
			return err
		}
		_, err = fmt.Println("Return accepted")
		if err != nil {
			return err
		}
	case GetReturnsCommand:
		_ = input.(GetReturnsCommandInput)
		data, err := h.useCase.GetReturns()
		if err != nil {
			return err
		}
		for _, order := range data {
			_, err := fmt.Println(order)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
