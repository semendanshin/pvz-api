package dummy

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecases"
)

var _ usecases.QueueProducer = &EventsSender{}

type EventsSender struct{}

func NewDummyEventsSender() *EventsSender {
	return &EventsSender{}
}

func (e EventsSender) SendEvent(_ context.Context, event domain.Event) error {
	fmt.Printf("Event sent: %v\n", event)
	return nil
}
