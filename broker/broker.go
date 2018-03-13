package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/Infra-Red/cf-postgresql-broker/config"
	"github.com/Infra-Red/cf-postgresql-broker/pgp"
	"github.com/pivotal-cf/brokerapi"
)

// New returns a new postgresql service broker instance.
func New(source string, services []brokerapi.Service, logger lager.Logger, env config.Specification) (brokerapi.ServiceBroker, error) {
	conn, err := pgp.New(source)
	if err != nil {
		return nil, err
	}
	return &postgresBroker{services: services, logger: logger, env: env, pgp: conn}, nil
}

type postgresBroker struct {
	services []brokerapi.Service
	logger   lager.Logger
	env      config.Specification
	pgp      *pgp.PGP
}

func (sb *postgresBroker) Services(context context.Context) ([]brokerapi.Service, error) {
	return sb.services, nil
}

func (sb *postgresBroker) Provision(context context.Context, instanceID string,
	details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	dbname, err := sb.pgp.CreateDB(context, instanceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}
	return brokerapi.ProvisionedServiceSpec{
		IsAsync:      false,
		DashboardURL: dbname,
	}, nil
}

func (sb *postgresBroker) Deprovision(context context.Context, instanceID string,
	details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, sb.pgp.DropDB(context, instanceID)
}

func (sb *postgresBroker) Bind(context context.Context, instanceID,
	bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	creds, err := sb.pgp.CreateUser(context, instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, err
	}
	return brokerapi.Binding{
		Credentials: creds,
	}, nil
}

func (sb *postgresBroker) Unbind(context context.Context, instanceID, bindingID string,
	details brokerapi.UnbindDetails) error {
	return sb.pgp.DropUser(context, instanceID, bindingID)
}

func (sb *postgresBroker) Update(context context.Context, instanceID string,
	details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

func (sb *postgresBroker) LastOperation(context context.Context, instanceID,
	operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
