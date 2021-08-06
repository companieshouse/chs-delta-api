package common

import "os"

const (
	TestRequestBody = `{"dummy" : "request"}`
	TestContextId   = "contextId"
	PostRestMethod = "POST"
)

func InitEnv() {
	_ = os.Setenv("BIND_ADDR", "bind_addr")
	_ = os.Setenv("KAFKA_BROKER_ADDR", "kafka_broker_addr,kafka_broker_addr")
	_ = os.Setenv("SCHEMA_REGISTRY_URL", "schema_registry_url")
	_ = os.Setenv("OFFICER_DELTA_TOPIC", "officer_delta_topic")
	_ = os.Setenv("OPEN_API_SPEC", "open_api_spec")
}

func DestroyEnv() {
	os.Clearenv()
}