package config

import "github.com/kelseyhightower/envconfig"

type Specification struct {
	BrokerPassword string `envconfig:"broker_password" required:"true"`
	BrokerUsername string `envconfig:"broker_username" required:"true"`
	CatalogPath    string `envconfig:"catalog_path" default:"./catalog.json"`
	LogLevel       string `envconfig:"log_level" default:"INFO"`
	PgSource       string `envconfig:"pg_source" required:"true"`
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
