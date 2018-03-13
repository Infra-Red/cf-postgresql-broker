package config

import "github.com/kelseyhightower/envconfig"

type Specification struct {
	BrokerUsername string `envconfig:"broker_username" required:"true"`
	BrokerPassword string `envconfig:"broker_password" required:"true"`
	PgSource       string `envconfig:"pg_source" required:"true"`
	LogLevel       string `envconfig:"log_level" default:"INFO"`
	Port           string `envconfig:"port" default:"8080"`
}

func LoadEnv() (Specification, error) {
	var env Specification
	err := envconfig.Process("", &env)
	if err != nil {
		return Specification{}, err
	}
	return env, nil
}
