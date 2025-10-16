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

# Create videos index
app_echo "Creating 'videos' index..."
curl -X PUT "$ELASTICSEARCH_HOST/videos" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/videos-mapping.json" || app_echo "Index 'videos' may already exist"

# Create video_manifests index
app_echo "Creating 'video_manifests' index..."
curl -X PUT "$ELASTICSEARCH_HOST/video_manifests" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/video_manifests-mapping.json" || app_echo "Index 'video_manifests' may already exist"

# Create video_variants index
app_echo "Creating 'video_variants' index..."
curl -X PUT "$ELASTICSEARCH_HOST/video_variants" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/video_variants-mapping.json" || app_echo "Index 'video_variants' may already exist"

# Create video_thumbnails index
app_echo "Creating 'video_thumbnails' index..."
curl -X PUT "$ELASTICSEARCH_HOST/video_thumbnails" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/video_thumbnails-mapping.json" || app_echo "Index 'video_thumbnails' may already exist"

# Create users index
app_echo "Creating 'users' index..."
curl -X PUT "$ELASTICSEARCH_HOST/users" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/users-mapping.json" || app_echo "Index 'users' may already exist"

# Create user_identities index
app_echo "Creating 'user_identities' index..."
curl -X PUT "$ELASTICSEARCH_HOST/user_identities" \
  -H 'Content-Type: application/json' \
  -d @"$SCRIPT_DIR/user_identities-mapping.json" || app_echo "Index 'user_identities' may already exist"

app_echo "All indices created successfully!"
