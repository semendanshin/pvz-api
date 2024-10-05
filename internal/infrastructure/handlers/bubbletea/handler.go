package bubbletea

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"homework/internal/abstractions"
)

// Handler is a handler for bubbletea
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
func (h *Handler) Run(ctx context.Context) error {
	models := make([]MyModel, 0)

	acceptOrderModel := newAcceptOrderModel(h.useCase)
	models = append(models, MyModel{
		Title: "Accept order",
		Model: acceptOrderModel,
	})

	returnOrderModel := newReturnOrderModel(h.useCase)
	models = append(models, MyModel{
		Title: "Return order",
		Model: returnOrderModel,
	})

	giveOrderToClientModel := newGiveOrderToClientModel(h.useCase)
	models = append(models, MyModel{
		Title: "Give order to client",
		Model: giveOrderToClientModel,
	})

	getOrdersModel := newGetOrdersModel(h.useCase, 10)
	models = append(models, MyModel{
		Title: "Get orders",
		Model: getOrdersModel,
	})

	acceptReturnModel := newAcceptReturnModel(h.useCase)
	models = append(models, MyModel{
		Title: "Accept return",
		Model: acceptReturnModel,
	})

	getReturnsModel := newGetReturnsModel(h.useCase, 10)
	models = append(models, MyModel{
		Title: "Get returns",
		Model: getReturnsModel,
	})

	p := tea.NewProgram(
		NewEntryPointModel(models),
		tea.WithMouseCellMotion(),
		// К сожалению, данный контекст используется только для остановки всего баблти приложения
		// В библиотеке не предусмотрена возможность использования этого контекста внутри моделей
		tea.WithContext(ctx),
	)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
