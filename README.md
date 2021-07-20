# chs-delta-api
Service to send delta events from CHIPS on the correct kafka topic

Environment Variables
-----------------

|  Variable                         |  Example                          |  Description                                       |
| --------------------------------- | --------------------------------- | -------------------------------------------------- |
| BIND_ADDR                         | 5010                              | Bind Address                                       |
| KAFKA_BROKER_ADDR                 | chs-kafka:9092                      | Kafka broker address (can be comma separated)    |
| SCHEMA_REGISTRY_URL               | http://chs-kafka:8081             | Schema registry URL                                |
| OFFICER_DELTA_TOPIC               | officers-delta                     | Kafka topic to write messages to                   |
| OPEN_API_SPEC                     | ./schema/officers                 | OpenAPI schema location                            |


## Running Locally with Docker CHS
Clone Docker CHS Development and follow the steps in the README.

Enable the `delta` module

All traffic will be handled via http://api.chs.local using the `/delta` endpoint. See API spec for available endpoints.

Development mode is available for this service in Docker CHS Development.

`./bin/chs-dev development enable chs-delta-api`

Swagger documentation is available for this service in the docker CHS development

Navigate to http://api.chs.local/api-docs/chs-delta-api/swagger-ui.html

## Running locally without Docker CHS
Pull image from private CH registry by running docker pull 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest command or run the following steps to build image locally:

1. `export SSH_PRIVATE_KEY_PASSPHRASE='[your SSH key passhprase goes here]'` (optional, set only if SSH key is passphrase protected)
2. `DOCKER_BUILDKIT=0 docker build --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" --build-arg SSH_PRIVATE_KEY_PASSPHRASE -t 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api .`
3. `docker run 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest`