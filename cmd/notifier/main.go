package main

import (
	"context"
	"fmt"
	"homework/internal/infrastructure/clients/queue/kafka"
	"homework/internal/infrastructure/saramawrapper/consumer_group"
	"log"
	"os"
	"sync"
)

func loadKafkaSettings() (brokers []string, topic string, group string) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	brokers = []string{broker}

	topic = os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "pvz.events-log"
	}

	group = os.Getenv("KAFKA_GROUP")
	if group == "" {
		group = "test-group"
	}

	return
}

func Run() error {
	wg := &sync.WaitGroup{}

	brokers, topic, group := loadKafkaSettings()
	log.Printf("brokers: %v, topic: %s, group: %s\n", brokers, topic, group)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := consumer_group.NewConsumerGroupHandler(kafka.NewEventsHandler())
	cg, err := consumer_group.NewConsumerGroup(
		brokers,
		group,
		[]string{topic},
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
