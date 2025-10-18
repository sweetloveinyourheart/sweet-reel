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

# Read connectors from configuration file
CONNECTORS_CONFIG="$SCRIPT_DIR/connectors.conf"

if [ -f "$CONNECTORS_CONFIG" ]; then
    app_echo "Reading connectors from $CONNECTORS_CONFIG"
    while IFS= read -r connector_name || [ -n "$connector_name" ]; do
        # Skip comments and empty lines
        [[ "$connector_name" =~ ^#.*$ ]] || [ -z "$connector_name" ] && continue

        # Trim whitespace
        connector_name=$(echo "$connector_name" | xargs)

        connector_file="$SCRIPT_DIR/${connector_name}.json"

        if [ -f "$connector_file" ]; then
            app_echo "Registering connector: $connector_name..."
            curl -X POST "$DEBEZIUM_HOST/connectors" \
              -H 'Content-Type: application/json' \
              -d @"$connector_file" || app_echo "Connector '$connector_name' may already exist"
        else
            app_echo "Warning: Connector file not found: $connector_file"
        fi
    done < "$CONNECTORS_CONFIG"
else
    app_echo "Warning: Connectors configuration file not found at $CONNECTORS_CONFIG"
fi

# Wait a bit for connectors to be fully initialized
sleep 10

# Note: Elasticsearch sink connectors require the Confluent Elasticsearch connector
# which may need to be installed separately in the Debezium image
app_echo "Note: Elasticsearch sink connectors require additional setup."
app_echo "You may need to register them manually after installing the Confluent connector."

app_echo "Checking connector status..."
curl -s "$DEBEZIUM_HOST/connectors" | jq '.'

app_echo "CDC connectors setup complete!"
