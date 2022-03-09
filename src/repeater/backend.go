package main

import (
	"owl/common/types"
)

type Backend interface {
	Write(data *types.TimeSeriesData) error
}
