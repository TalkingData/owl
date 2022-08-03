package backend

import (
	"fmt"
	"owl/dto"
	"owl/repeater/conf"
)

// Backend interface
type Backend interface {
	Write(data *dto.TsData) error
}

// NewBackend
func NewBackend(conf *conf.Conf) (Backend, error) {
	switch conf.Backend {
	case "opentsdb":
		return newOpentsdbBackend(conf.OpentsdbAddress)
	case "kairosdb":
		return newOpentsdbBackend(conf.KairosdbAddress)
	case "kairosdb-rest":
		return newKairosdbRestBackend(conf.KairosdbRestAddress)
	case "kafka":
		return newKafkaBackend(conf.KafkaAddresses, conf.KafkaTopic)
	}

	return nil, fmt.Errorf("unsupported backend %s", conf.Backend)
}
