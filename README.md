# chs-delta-api
Service to send delta events from CHIPS on the correct kafka topic

Environment Variables
-----------------

| Variable                          | Example                  | Description                                           |  Required       | Default value |
|-----------------------------------|--------------------------|-------------------------------------------------------| --------------- | ------------- |
| BIND_ADDR                         | 5010                     | Bind Address / application port                       | YES             |               |
| KAFKA_BROKER_ADDR                 | chs-kafka:9092           | Kafka broker address (can be comma separated)         | YES             |               |
| SCHEMA_REGISTRY_URL               | http://chs-kafka:8081    | Schema registry URL                                   | YES             |               |
| OFFICER_DELTA_TOPIC               | officers-delta           | Officer Delta Kafka topic to write messages to        | YES             |               |
| INSOLVENCY_DELTA_TOPIC            | insolvency-delta         | Insolvency Delta Kafka topic to write messages to     | YES             |               |
| CHARGES_DELTA_TOPIC               | charges-delta            | Charges Delta Kafka topic to write messages to        | YES             |               |
| DISQUALIFIED_OFFICERS_DELTA_TOPIC | disqualification-delta   | Disqualified Delta Kafka topic to write messages to   | YES             |               |
| COMPANY_DELTA_TOPIC               | company-profile-delta    | Company Delta Kafka topic to write messages to        | YES             |               |
| EXEMPTION_DELTA_TOPIC             | company-exemptions-delta | Exemption Delta Kafka topic to write messages to      | YES             |               |
| PSC_STATEMENT_DELTA_TOPIC         | psc-statement-delta      | Psc Statement Delta Kafka topic to write messages to  | YES             |               |
| FILING_HISTORY_DELTA_TOPIC        | filing-history-delta     | Filing History Delta Kafka topic to write messages to | YES             |               |
| DOCUMENT_STORE_DELTA_TOPIC        | document-store-delta     | Document Store Delta Kafka topic to write messages to | YES             |               |
| REGISTERS_DELTA_TOPIC             | registers-delta          | Registers Delta Kafka topic to write messages to      | YES             |               |
| ACSP_PROFILE_DELTA_TOPIC          | acsp-profile-delta       | ACSP Profile Delta Kafka topic to write messages to   | YES             |               |
| OPEN_API_SPEC                     | ./apispec/api-spec.yml   | OpenAPI schema location                               | YES             |               |
| LOG_LEVEL                         | trace                    | The level at which the logger prints                  | NO              | info          |

## Running Locally with Docker CHS
Clone [Docker CHS Development](https://github.com/companieshouse/docker-chs-development) and follow the steps in the README.

### Using the delta module
Enable the `delta` module  by running the command `./bin/chs-dev modules enable delta` in chs-docker-development

All traffic will be handled via http://api.chs.local:4001 using the `/delta` endpoint. See API spec for available endpoints.
### Running the service on its own
To run the service on its own use the command `./bin/chs-dev services enable chs-delta-api` in chs-docker-development

### Running the service in Development Mode
Development mode is available for this service in Docker CHS Development.

`./bin/chs-dev development enable chs-delta-api`

## Running locally in Docker, without Docker CHS
Export the required environment variables to ensure the service can start up correctly.

Pull image from private CH registry by running `docker pull 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest` 
command or run the following steps to build image locally:

1. `export SSH_PRIVATE_KEY_PASSPHRASE='[your SSH key passhprase goes here]'` (optional, set only if SSH key is passphrase protected)
2. `DOCKER_BUILDKIT=0 docker build --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" --build-arg SSH_PRIVATE_KEY_PASSPHRASE -t 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api .`
3. `docker run 169942020521.dkr.ecr.eu-west-1.amazonaws.com/local/chs-delta-api:latest`

Local testing
=============

## Getting Started

To build the service:
 1. Clone the repository into your GOPATH under `src/github.com/companieshouse`
 2. Build the docker image using the image definiton from the docker-chs-development repo using:
 ```shell
 make docker-image
 ```
 3. Run the service in the chs-env by enabling the filing-notification-sender module:
 ```shell
  chs-dev servies enable authentication-service chs-delta-api officer-delta-processor company-appointments-api-ch-gov-uk company-metrics-api company-metrics-consumer

  chs-dev modules enable streaming

  chs-dev up
   ```
4. When run in the background logs can be found at
  chs-dev logs -f chs-delta-api

## Healthcheck
This service implements a `healthcheck` endpoint. Using POSTMAN call the `/delta/healthcheck` GET endpoint to assert 
the service is running correctly.

## Documentation
All documentation can be found in the `/docs` folder at the root of this project's directory.

Note : Make sure you are logged into AWS ECR before running any of the above commands
`aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 169942020521.dkr.ecr.eu-west-1.amazonaws.com
`
