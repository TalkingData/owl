package backend

import (
	"owl/common/types"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaBackend struct {
	Producer sarama.AsyncProducer
	Topic    string
}

func NewKafkaBackend(brokers []string, topic string) (*KafkaBackend, error) {
	producer, err := sarama.NewAsyncProducer(brokers, nil)
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
		Key:   nil,
		Value: sarama.StringEncoder(string(data.Encode())),
	}
	select {
	case backend.Producer.Input() <- message:
		return nil
	case err := <-backend.Producer.Errors():
		return err
	}
}
