package backend

import (
	"gopkg.in/Shopify/sarama.v1"
	"owl/dto"
)

// kafkaBackend struct
type kafkaBackend struct {
	Producer sarama.SyncProducer
	Topic    string
}

// newKafkaBackend
func newKafkaBackend(addresses []string, topic string) (Backend, error) {
	producer, err := sarama.NewSyncProducer(addresses, nil)
	if err != nil {
		return nil, err
	}
	bEnd := &kafkaBackend{
		Producer: producer,
		Topic:    topic,
	}
	return bEnd, nil
}

// Write
func (kfk *kafkaBackend) Write(data *dto.TsData) error {
	message := &sarama.ProducerMessage{
		Topic: kfk.Topic,
		Value: sarama.StringEncoder(data.Encode()),
	}
	_, _, err := kfk.Producer.SendMessage(message)
	return err
}
