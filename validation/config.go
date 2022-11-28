package validation

import "src.sqlkite.com/utils/log"

type Config struct {
	PoolSize  uint16 `json:"pool_size"`
	MaxErrors uint16 `json:"max_errors"`
}

func Configure(config Config) error {
	poolSize := config.PoolSize
	if poolSize == 0 {
		poolSize = 100
	}

	maxErrors := config.MaxErrors
	if maxErrors == 0 {
		maxErrors = 20
	}

	globalPool = NewPool(poolSize, maxErrors)
	log.Info("validation_config").
		Int("pool_size", int(poolSize)).
		Int("max_errors", int(maxErrors)).
		Log()

	return nil
}
