package logger

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/Infra-Red/cf-postgresql-broker/config"
)

var logLevels = map[string]lager.LogLevel{
	"DEBUG": lager.DEBUG,
	"INFO":  lager.INFO,
	"ERROR": lager.ERROR,
	"FATAL": lager.FATAL,
}

func NewLogger(component string, env config.Specification) (lager.Logger, error) {
	logger := lager.NewLogger(component)
	logLevel, ok := logLevels[env.LogLevel]
	if !ok {
		return nil, fmt.Errorf("Unknown log level %s. Available log levels are: DEBUG, INFO, ERROR, and FATAL", env.LogLevel)
	}
	fmt.Printf("Using log level %s\n", env.LogLevel)
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, logLevel))
	return logger, nil
}
