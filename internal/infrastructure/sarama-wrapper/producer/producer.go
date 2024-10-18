package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"homework/internal/infrastructure/sarama-wrapper"
)

func NewSyncSaramaProducer(conf sarama_wrapper.Config, opts ...Option) (sarama.SyncProducer, error) {
	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}
