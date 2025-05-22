#!/bin/bash
#
#
# Start script for chs-delta-api
PORT="5010"

# Read brokers and topics from environment and split on comma
IFS=',' read -ra BROKERS <<< "${KAFKA_BROKER_ADDR}"

# Ensure we only populate the broker address and topic via application arguments
unset KAFKA_BROKER_ADDR

exec ./chs-delta-api $(for broker in "${BROKERS[@]}"; do echo -n "-broker-addr=${broker} "; done)