# chs-delta-api
Service to send delta events from CHIPS on the correct kafka topic

## Requirements

In order to build the service locally you will need the following:

- [Go](https://go.dev/doc/install)
- [Git](https://git-scm.com/downloads)

These should be installed/updated using the CH dev env scripts.

Note the following:

- Version 1.22 of Go should be installed in your CHS development environment
- Check settings for the GOPRIVATE, GOPROXY and GONOPROXY environment variables using the command `go env`. Working values are:

```shell
GOPRIVATE='github.com/companieshouse/*'
GOPROXY='https://proxy.golang.org,direct'
GONOPROXY='github.com/pierrec/*,gopkg.in/jcmturner/*'
```

- If certificate download errors are encountered when building or running tests, try disconnecting from the VPN and running `chproxyoff` before reissuing the command (reconnecting and reenabling the proxy will need to be done afterwards)
- Running `make test` from the terminal will build the project and run the unit tests


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

After making local code changes, build a new 'development' Docker image using:

```shell
go mod tidy
OOS=linux GOARCH=amd64 GOOS=linux go build
```

A new image named `chs-delta-api` should be created in the project folder and Docker CHS must then be restarted for the changes to take effect.

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
This service implements a `healthcheck` endpoint. Using POSTMAN call the `/chs-delta-api/healthcheck` GET endpoint to assert 
the service is running correctly.

## Documentation
All documentation can be found in the `/docs` folder at the root of this project's directory.

Note : Make sure you are logged into AWS ECR before running any of the above commands
`aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 169942020521.dkr.ecr.eu-west-1.amazonaws.com
`
## Terraform ECS

### What does this code do?

The code present in this repository is used to define and deploy a dockerised container in AWS ECS.
This is done by calling a [module](https://github.com/companieshouse/terraform-modules/tree/main/aws/ecs) from terraform-modules. Application specific attributes are injected and the service is then deployed using Terraform via the CICD platform 'Concourse'.


Application specific attributes | Value                                | Description
:---------|:-----------------------------------------------------------------------------|:-----------
**ECS Cluster**        | data-sync                                     | ECS cluster (stack) the service belongs to
**Load balancer**      |{env}-chs-internalapi                                           | The load balancer that sits in front of the service
**Concourse pipeline**     |[Pipeline link](https://ci-platform.companieshouse.gov.uk/teams/team-development/pipelines/chs-delta-api) <br> [Pipeline code](https://github.com/companieshouse/ci-pipelines/blob/master/pipelines/ssplatform/team-development/chs-delta-api)                                  | Concourse pipeline link in shared services


### Contributing
- Please refer to the [ECS Development and Infrastructure Documentation](https://companieshouse.atlassian.net/wiki/spaces/DEVOPS/pages/4390649858/Copy+of+ECS+Development+and+Infrastructure+Documentation+Updated) for detailed information on the infrastructure being deployed.

### Testing
- Ensure the terraform runner local plan executes without issues. For information on terraform runners please see the [Terraform Runner Quickstart guide](https://companieshouse.atlassian.net/wiki/spaces/DEVOPS/pages/1694236886/Terraform+Runner+Quickstart).
- If you encounter any issues or have questions, reach out to the team on the **#platform** slack channel.

### Vault Configuration Updates
- Any secrets required for this service will be stored in Vault. For any updates to the Vault configuration, please consult with the **#platform** team and submit a workflow request.

### Useful Links
- [ECS service config dev repository](https://github.com/companieshouse/ecs-service-configs-dev)
- [ECS service config production repository](https://github.com/companieshouse/ecs-service-configs-production)
