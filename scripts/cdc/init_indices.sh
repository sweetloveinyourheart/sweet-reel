#!/bin/bash
. ./scripts/util.sh

set -e

ELASTICSEARCH_HOST=${ELASTICSEARCH_HOST:-"http://localhost:9200"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/elasticsearch"

app_echo "Waiting for Elasticsearch to be ready..."
until curl -s "$ELASTICSEARCH_HOST/_cluster/health" > /dev/null; do
  app_echo "Elasticsearch is not ready yet. Retrying in 5 seconds..."
  sleep 5
done

app_echo "Elasticsearch is ready. Creating indices..."

# Read indices from configuration file
INDICES_CONFIG="$SCRIPT_DIR/indices.conf"

if [ -f "$INDICES_CONFIG" ]; then
    app_echo "Reading indices from $INDICES_CONFIG"
    while IFS=':' read -r index_name mapping_name || [ -n "$index_name" ]; do
        # Skip comments and empty lines
        [[ "$index_name" =~ ^#.*$ ]] || [ -z "$index_name" ] && continue

        # Trim whitespace
        index_name=$(echo "$index_name" | xargs)
        mapping_name=$(echo "$mapping_name" | xargs)

        # Use index_name as mapping_name if not specified
        mapping_name=${mapping_name:-$index_name}

        mapping_file="$SCRIPT_DIR/${mapping_name}-mapping.json"

        if [ -f "$mapping_file" ]; then
            app_echo "Creating '$index_name' index..."
            curl -X PUT "$ELASTICSEARCH_HOST/$index_name" \
              -H 'Content-Type: application/json' \
              -d @"$mapping_file" || app_echo "Index '$index_name' may already exist"
        else
            app_echo "Warning: Mapping file not found: $mapping_file"
        fi
    done < "$INDICES_CONFIG"
else
    app_echo "Warning: Indices configuration file not found at $INDICES_CONFIG"
fi

app_echo "All indices created successfully!"
