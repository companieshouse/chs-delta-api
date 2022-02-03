# Adding a new delta

## Overview

To add a new delta to the chs-delta-api you need to complete the following steps:
1. Create the openAPI 3 spec for the new delta inside of the `/apispec` directory and add its route to the main
`api-spec.yml` file. Also, remembering to add unit tests inside of the `/validation/schema_testing` directory to cover your 
new specs functionality (for more details on schema testing, see `unit-testing-a-new-schema` documentation in the `/docs` directory).
2. Add the new delta to the `Config` struct inside of the `/config/Config.go` file. Also, update the `/config/config_test.go` 
test to account for your new delta.
3. Add your new route to the `register.go` file inside of the `/handlers` directory, using the generic DeltaHandler provided.

## 1. Creating the OpenAPI spec
Inside of the `/apispec` directory create a new yml file (e.g. `example-delta-spec.yml`). Inside of the new spec file create
your delta spec.

Finally, associate the new delta-spec.yml file with a route by adding it to the `api-spec.yml` file under the paths section.
```yaml
paths:
  /delta/example-delta:
    $ref: 'example-delta-spec.yml'
```

## 2. Adding the new delta to config
Inside of the `/config/Config.go` file add to the Config struct a new string environment variable to hold your delta's topic name.
```go
ExampleDeltaTopic string `env:"EXAMPLE_DELTA_TOPIC" flag:"example-delta-topic" flagDesc:"Topic for example delta"`
```

## 3. Exposing your new delta endpoint
You need to register your new route within the `register.go` file.
```go
appRouter.HandleFunc("/delta/example-delta", NewDeltaHandler(kSvc, h, chv, cfg, false, cfg.ExampleDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("example-delta")
```

You may also want to optionally add a validation endpoint for your new delta which only handles validation and doesn't send 
the request to Kafka. to do this, switch the boolean parameter (doValidationOnly) on the DeltaHandler constructor to true as seen in the snippet below.
```go
appRouter.HandleFunc("/delta/example-delta/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, cfg.ExampleDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("example-delta-validate")
```

Finally, you'll need to update the register.go `TestUnitRegister` unit test to cover your changes.

## 4. Updating docker compose to specify your kafka topic
In order to run this you need to have added a new environment variable in the respective docker compose file located in the [docker-chs-development repo](https://github.com/companieshouse/docker-chs-development). For deltas this is `services/modules/delta/chs-delta-api.docker-compose.yaml`. 

Add your kafka topic name under `services.chs-delta-api.environment` like so:

```yaml
services:
  chs-delta-api:
    ...
    environment:
      ...
      - EXAMPLE_DELTA_TOPIC=example-delta
      ...
```

## Final notes
All other services on the chs-delta-api are generic and will not require any changes when adding a new delta, including 
their unit tests.
