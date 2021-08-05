# Adding a new delta

## Overview

To add a new delta to the chs-delta-api you need to complete the following steps:
1. Create the API 3 spec for the new delta inside of the `/apispec` directory and add its route to the main
`api-spec.yml` file. Also remembering to add unit tests inside of the `/validation/schema_testing` directory to cover your 
new specs functionality (for more details on schema testing, see the schema testing documentation in the `/docs` directory).
2. Add the new delta to the `Config` struct inside of the `/config/Config.go` file. Also update the `/config/config_test.go` 
test to account for your new delta.
3. Create a handler inside of the `/handlers` directory to handle the route associated with the new delta along with 
associated unit tests.
4. Add your new route to the `restiger.go` file inside of the `/handlers` directory.

## 1. Creating the API spec
Inside of the `/apispec` directory create a new yml file (e.g. `example-delta-spec.yml`). Inside of the new spec file create
your delta spec.

Finally, associate the new delta-spec.yml file with a route by adding it to the `api-sepc.yml` file under the paths section.
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

## 3. Creating the handler
Inside of the `/handlers` directory create a new Go file and name it appropriately (e.g. `exampleHandler.go`). 
Inside of the new handler file create a struct which implements the Handler interface. Use the following template
as a guide.

- The KafkaService will be needed to send your data onto a given Kafka topic.
- The helper struct will be needed to extract the contextId from your request and to marshal your request into a string 
ready for sending to Kafka.
- The CHValidator will be used to validate your request against your newly created api-spec.
- The config struct will be used to load all needed configs.

```go
type ExampleHandler struct {
	kSvc services.KafkaService
	h    helpers.Helper
	chv  validation.CHValidator
	cfg  *config.Config
}

func NewExampleDeltaHandler(kSvc services.KafkaService, h helpers.Helper, chv validation.CHValidator, cfg *config.Config) *ExampleDeltaHandler {
	return &ExampleDeltaHandler{kSvc: kSvc, h: h, chv: chv, cfg: cfg}
}


func (eh *ExampleDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextId := eh.h.GetRequestIdFromHeader(r)
    
    log.InfoC(contextId, fmt.Sprintf("Using the open api spec: "), log.Data{config.OpenApiSpecKey: eh.cfg.OpenApiSpec})

    // Validate against the open API 3 spec before progressing any further.
    errValidation, err := eh.chv.ValidateRequestAgainstOpenApiSpec(r, eh.cfg.OpenApiSpec, contextId)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while trying to validate request"})
        return
    } else if errValidation != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, err = w.Write(errValidation)
        if err != nil {
            log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while trying to write response"})
        }

        return
    }

    // Get request body and marshal into a string, ready for publishing.
    data, err := eh.h.GetDataFromRequest(r, contextId)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Send data string to Kafka service for publishing.
    if err := eh.kSvc.SendMessage(eh.cfg.ExampleDeltaTopic, data, contextId); err != nil {
        log.ErrorC(contextId, err, log.Data{config.TopicKey: eh.cfg.ExampleDeltaTopic, config.MessageKey: "error sending the message to the given kafka topic"})
        w.WriteHeader(http.StatusInternalServerError)

        return
    }

    w.WriteHeader(http.StatusOK)
}
```

## 4. Exposing your new delta endpoint
The final step of adding a new delta is to register your new route within the `register.go` file.
```go
appRouter.HandleFunc("/delta/example-delta", NewExampleDeltaHandler(kSvc, h, chv, cfg).ServeHTTP).Methods(http.MethodPost).Name("example-delta")
```
and update the register.go `TestRegister` unit test to cover your changes.

## Final notes
All other services on the chs-delta-api are generic and will not require any changes when adding a new delta, including 
their unit tests.