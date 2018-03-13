build:
	go build -o "postgresql-broker" github.com/Infra-Red/cf-postgresql-broker/cmd/postgresql-broker

push:
	cf push

register:
	cf create-service-broker shared-postgresql-broker admin admin http://postgresql-broker.bosh-lite.com
	cf enable-service-access postgresql