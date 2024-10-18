package main

import (
	"context"
	"fmt"
	"homework/internal/infrastructure/clients/queue/kafka"
	"homework/internal/infrastructure/sarama-wrapper/consumer_group"
	"log"
	"os"
	"sync"
)

func Run() error {
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := consumer_group.NewConsumerGroupHandler(kafka.NewEventsHandler())
	cg, err := consumer_group.NewConsumerGroup(
		[]string{"localhost:9092"},
		"test-group",
		[]string{"pvz.events-log"},
		handler,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cg.Close()

	cg.Run(ctx, wg)

	wg.Wait()

	return nil
}

func main() {
	if err := Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
