package config

//BindAddr is the port on which the app listen on
const BindAddr = "bind_address"

//BrokerAddr is the address of the Kafka broker
const BrokerAddr = "broker_address"

//SchemaRegistryURL is the url where chs-delta schema is registered
const SchemaRegistryURL = "schema_registry_url"

//OfficerDeltaTopic is the name of the Kafka topic
const OfficerDeltaTopic = "officer_delta_topic"

//OpenApiSpec is the location where open api schema is stored
const OpenApiSpec = "open_api_spec"

//SchemaAbsolutePath is the key to get the absolute path to the schema
const SchemaAbsolutePath = "schema_absolute_path"
