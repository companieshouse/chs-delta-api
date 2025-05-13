#!/bin/bash
#
# Start script for chs-delta-api
PORT="5010"

exec ./chs-delta-api "-bind-addr=:${PORT}"