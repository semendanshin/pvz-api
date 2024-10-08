package stdin

import (
	"context"
	"fmt"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"strconv"
	"strings"
	"time"
)

const (
	AcceptDeliveryCommand      Command = "accept-delivery"
	AcceptReturnCommand        Command = "accept-return"
	GetOrdersCommand           Command = "get-orders"
	GetReturnsCommand          Command = "get-returns"
	GiveOrderToClientCommand   Command = "give-orders"
	ReturnOrderDeliveryCommand Command = "return-delivery"
)

type Handler struct {
	useCase abstractions.IPVZOrderUseCase
	srv     *Server
}

func NewHandler(useCase abstractions.IPVZOrderUseCase, numOfWorkers int) *Handler {
	return &Handler{
		useCase: useCase,
		srv: NewServer(
			numOfWorkers,
		),
	}
}

func (h *Handler) Run(ctx context.Context) error {
	h.srv.AddHandler(AcceptDeliveryCommand, h.AcceptDeliveryHandler)
	h.srv.AddHandler(AcceptReturnCommand, h.AcceptReturnHandler)
	h.srv.AddHandler(GetOrdersCommand, h.GetOrdersHandler)
	h.srv.AddHandler(GetReturnsCommand, h.GetReturnsHandler)
	h.srv.AddHandler(GiveOrderToClientCommand, h.GiveOrderToClientHandler)
	h.srv.AddHandler(ReturnOrderDeliveryCommand, h.ReturnOrderDeliveryHandler)

	return h.srv.Run(ctx)
}

func (h *Handler) Stop() {
	h.srv.Stop()
}

func (h *Handler) AcceptDeliveryHandler(ctx context.Context, args []string) (string, error) {
	usage := "<order_id> <recipient_id> <storage_time: 1h30m> <cost> <weight> <packaging> ?<additional_film: bool>"

	if len(args) < 6 || len(args) > 7 {
		return "", fmt.Errorf("invalid number of arguments, expected 6 or 7, got %d. Usage: %s", len(args), usage)
	}

	var input struct {
		OrderID        string
		RecipientID    string
		StorageTime    time.Duration
		Cost           int
		Weight         int
		Packaging      domain.PackagingType
		AdditionalFilm bool
	}
	{
		var err error

		input.OrderID = args[0]

		input.RecipientID = args[1]

		input.StorageTime, err = time.ParseDuration(args[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse storage time: %w", err)
		}

		if input.StorageTime < 0 {
			return "", fmt.Errorf("storage time is negative")
		}

		input.Cost, err = strconv.Atoi(args[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse cost: %w", err)
		}

		input.Weight, err = strconv.Atoi(args[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse weight: %w", err)
		}

		input.Packaging, err = domain.NewPackagingType(args[5])
		if err != nil {
			return "", fmt.Errorf("failed to parse packaging: %w", err)
		}

		if len(args) == 7 {
			input.AdditionalFilm, err = strconv.ParseBool(args[6])
			if err != nil {
				return "", fmt.Errorf("failed to parse additional film: %w", err)
			}
		}
	}

	err := h.useCase.AcceptOrderDelivery(
		ctx,
		input.OrderID,
		input.RecipientID,
		input.StorageTime,
		input.Cost,
		input.Weight,
		input.Packaging,
		input.AdditionalFilm,
	)
	if err != nil {
		return "", err
	}

	return "Delivery accepted", nil
}

func (h *Handler) AcceptReturnHandler(ctx context.Context, args []string) (string, error) {
	usage := "<recipient_id> <order_id>"

	if len(args) != 2 {
		return "", fmt.Errorf("invalid number of arguments, expected 2, got %d. Usage: %s", len(args), usage)
	}

	recipientID := args[0]
	orderID := args[1]

	err := h.useCase.AcceptReturn(ctx, recipientID, orderID)
	if err != nil {
		return "", err
	}

	return "Return accepted", nil
}

func (h *Handler) GetOrdersHandler(ctx context.Context, args []string) (string, error) {
	// Need to do smth with SamePVZ, LastN, Cursor and Limit options
	usage := "<user_id>"

	if len(args) != 1 {
		return "", fmt.Errorf("invalid number of arguments, expected 1, got %d. Usage: %s", len(args), usage)
	}

	userID := args[0]

	orders, err := h.useCase.GetOrders(ctx, userID)
	if err != nil {
		return "", err
	}

	strOrders := make([]string, len(orders))
	for i, order := range orders {
		strOrders[i] = fmt.Sprintf("%s %s %s %d %d %s %t",
			order.OrderID,
			order.RecipientID,
			order.PVZID,
			order.Cost,
			order.Weight,
			order.Packaging,
			order.AdditionalFilm,
		)
	}
	return strings.Join(strOrders, "\n"), nil
}

func (h *Handler) GetReturnsHandler(ctx context.Context, _ []string) (string, error) {
	// Need to do smth with Limit and Offset options
	orders, err := h.useCase.GetReturns(ctx)
	if err != nil {
		return "", err
	}

	strOrders := make([]string, len(orders))
	for i, order := range orders {
		strOrders[i] = fmt.Sprintf("%s %s %s %d %d %s %t",
			order.OrderID,
			order.RecipientID,
			order.PVZID,
			order.Cost,
			order.Weight,
			order.Packaging,
			order.AdditionalFilm,
		)
	}

	return strings.Join(strOrders, "\n"), nil
}

func (h *Handler) GiveOrderToClientHandler(ctx context.Context, args []string) (string, error) {
	usage := "<order_id1> <order_id2> ..."

	if len(args) < 1 {
		return "", fmt.Errorf("invalid number of arguments, expected at least 1, got %d. Usage: %s", len(args), usage)
	}

	err := h.useCase.GiveOrderToClient(ctx, args)
	if err != nil {
		return "", err
	}

	return "Orders given to client", nil
}

func (h *Handler) ReturnOrderDeliveryHandler(ctx context.Context, args []string) (string, error) {
	usage := "<order_id>"

	if len(args) != 1 {
		return "", fmt.Errorf("invalid number of arguments, expected 1, got %d. Usage: %s", len(args), usage)
	}

	orderID := args[0]

	err := h.useCase.ReturnOrderDelivery(ctx, orderID)
	if err != nil {
		return "", err
	}

	return "Delivery returned", nil
}
