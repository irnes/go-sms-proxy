#!/bin/bash

originator="IRNES-API"
recipient=$1
message=$2

API="http://127.0.0.1:8080/messages"
PAYLOAD="{\"recipient\": \"$recipient\", \"originator\": \"$originator\", \"message\": \"$message\"}"

echo $PAYLOAD

curl -i -H 'Content-Type: application/json' -d "$PAYLOAD" $API

echo ""
