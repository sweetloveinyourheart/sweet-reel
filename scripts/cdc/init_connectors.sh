#!/bin/bash
. ./scripts/util.sh

set -e

DEBEZIUM_HOST=${DEBEZIUM_HOST:-"http://localhost:8083"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/debezium"

app_echo "Waiting for Debezium Connect to be ready..."
until curl -s "$DEBEZIUM_HOST/" > /dev/null; do
  app_echo "Debezium Connect is not ready yet. Retrying in 5 seconds..."
  sleep 5
done

app_echo "Debezium Connect is ready. Registering connectors..."

# Register video-management connector
app_echo "Registering video-management connector..."
curl -X POST "$DEBEZIUM_HOST/connectors" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/video-management-connector.json" || app_echo "Connector 'video-management-connector' may already exist"

# Register user connector
app_echo "Registering user connector..."
curl -X POST "$DEBEZIUM_HOST/connectors" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/user-connector.json" || app_echo "Connector 'user-connector' may already exist"

# Wait a bit for connectors to be fully initialized
sleep 10

# Note: Elasticsearch sink connectors require the Confluent Elasticsearch connector
# which may need to be installed separately in the Debezium image
app_echo "Note: Elasticsearch sink connectors require additional setup."
app_echo "You may need to register them manually after installing the Confluent connector."

app_echo "Checking connector status..."
curl -s "$DEBEZIUM_HOST/connectors" | jq '.'

app_echo "CDC connectors setup complete!"
