# PostgreSQL Broker
[![Maintainability](https://api.codeclimate.com/v1/badges/8f1ee2ab2e2b21822e09/maintainability)](https://codeclimate.com/github/Infra-Red/cf-postgresql-broker/maintainability) [![Build Status](https://semaphoreci.com/api/v1/infra-red/cf-postgresql-broker/branches/master/badge.svg)](https://semaphoreci.com/infra-red/cf-postgresql-broker)

This service broker shares a PostgreSQL amongst many users via the Open Service Broker API.

## Configuration

There are important environment variables that should be overriden inside the manifest.yml file:

* `BROKER_USERNAME` and `BROKER_PASSWORD` are required to setup [HTTP Basic Auth](https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#authentication) to the broker API.
* `PG_SOURCE` - URL used for connection to a PostgreSQL database instance. Can be customized according to [this library](https://godoc.org/github.com/lib/pq).

## Deployment

1. Clone this repository, and `cd` into it.
1. Target the space you want to deploy the broker to.

    ```bash
    $ cf target -o <org> -s <space>
    ```
1. The configuration is entirely read from environment variables. Edit the manifest.yml according to your environment.
1. Deploy the broker as an application.

    ```bash
    $ cf push
    ```
1. [Register a Broker](http://docs.cloudfoundry.org/services/managing-service-brokers.html#register-broker).

    ```bash
    $ cf create-service-broker shared-postgresql-broker $BROKER_USERNAME $BROKER_PASSWORD $APP-URL
    ```
1. [Enable Access to Service Plans](http://docs.cloudfoundry.org/services/access-control.html#enable-access).

    ```bash
    $ cf enable-service-access a.postgresql
    ```

## Using

All credentials stored in the `VCAP_SERVICES` environment variable with the JSON key `a.postgresql`.

1. [Create a Service Instance](https://docs.cloudfoundry.org/devguide/services/managing-services.html#create).

    ```bash
    $ cf create-service a.postgresql standard my-postgresql
    ```

1. [Bind a Service Instance](https://docs.cloudfoundry.org/devguide/services/managing-services.html#bind).

    ```bash
    $ cf bind-service my-app my-postgresql
    ```

## Demo

[![asciicast](https://asciinema.org/a/q8GV3pcDWzSAtskdjtQ2v2Nl7.png)](https://asciinema.org/a/q8GV3pcDWzSAtskdjtQ2v2Nl7)

## License

```
Copyright 2017 Andrei Krasnitski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```