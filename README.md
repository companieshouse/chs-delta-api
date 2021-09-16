# chs-delta-api
Service to send delta events from CHIPS on the correct kafka topic

Environment Variables
-----------------

|  Variable                         |  Example                          |  Description                                       |  Required       | Default value |        
| --------------------------------- | --------------------------------- | -------------------------------------------------- | --------------- | ------------- |
| BIND_ADDR                         | 5010                              | Bind Address / application port                    | YES             |               |
| KAFKA_BROKER_ADDR                 | chs-kafka:9092                    | Kafka broker address (can be comma separated)      | YES             |               |
| SCHEMA_REGISTRY_URL               | http://chs-kafka:8081             | Schema registry URL                                | YES             |               |
| OFFICER_DELTA_TOPIC               | officers-delta                    | Officer Delta Kafka topic to write messages to     | YES             |               |
| OPEN_API_SPEC                     | ./apispec/apispec.yml             | OpenAPI schema location                            | YES             |               |
| LOG_LEVEL                         | trace                             | The level at which the logger prints               | NO              | info          |

## Running Locally with Docker CHS
Clone [Docker CHS Development](https://github.com/companieshouse/docker-chs-development) and follow the steps in the README.

Enable the `delta` module

All traffic will be handled via http://api.chs.local:4001 using the `/delta` endpoint. See API spec for available endpoints.

Development mode is available for this service in Docker CHS Development.

`./bin/chs-dev development enable chs-delta-api`

## Running locally in Docker, without Docker CHS
Export the required environment variables to ensure the service can start up correctly.

Pull image from private CH registry by running `docker pull 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest` 
command or run the following steps to build image locally:

1. `export SSH_PRIVATE_KEY_PASSPHRASE='[your SSH key passhprase goes here]'` (optional, set only if SSH key is passphrase protected)
2. `DOCKER_BUILDKIT=0 docker build --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" --build-arg SSH_PRIVATE_KEY_PASSPHRASE -t 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api .`
3. `docker run 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest`

## Healthcheck
This service implements a `healthcheck` endpoint. Using POSTMAN call the `/delta/healthcheck` GET endpoint to assert 
the service is running correctly.

## Documentation
All documentation can be found in the `/docs` folder at the root of this project's directory.