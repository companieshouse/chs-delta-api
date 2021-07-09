package models

// ChsDelta is a struct that replicates the structure of the chs-delta avro.
type ChsDelta struct {
	CreatedAt string `avro:"created_at"`
	DeltaAt   int32 `avro:"delta_at"`
	Data      string `avro:"data"`
	Attempt   int32  `avro:"attempt"`
	ContextId string `avro:"context_id"`
}