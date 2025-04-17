#!/bin/bash
#
# Start script for chs-monitor-api
PORT="5010"

exec ./chs-kafka-api "-bind-addr=:${PORT}"