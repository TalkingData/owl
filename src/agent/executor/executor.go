package executor

import "owl/common/logger"

type Executor struct {
	logger *logger.Logger
}

func NewExecutor(lg *logger.Logger) *Executor {
	return &Executor{logger: lg}
}
