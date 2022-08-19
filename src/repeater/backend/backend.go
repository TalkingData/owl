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
	case "kairosdb":
		return newKairosDbBackend(conf.KairosDbAddress)
	}

	return nil, fmt.Errorf("unsupported backend %s", conf.Backend)
}
