package backend

import (
	"owl/common/types"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaBackend struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewKafkaBackend(brokers []string, topic string) (*KafkaBackend, error) {
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		return nil, err
	}
	backend := &KafkaBackend{
		Producer: producer,
		Topic:    topic,
	}
	return backend, nil
}

func (backend *KafkaBackend) Write(data *types.TimeSeriesData) error {
	message := &sarama.ProducerMessage{
		Topic: backend.Topic,
		Value: sarama.StringEncoder(string(data.Encode())),
	}
	_, _, err := backend.Producer.SendMessage(message)
	return err
}
