#!/bin/bash

set -euo pipefail

(minio server /data --console-address ":9001" &)

# Wait for MinIO to be ready
echo "Waiting for MinIO to start..."
until curl -s http://localhost:9000/minio/health/ready > /dev/null; do
  sleep 1
done
echo "MinIO is ready."

# Configure mc alias
echo "Configuring mc alias..."
mc alias set local http://localhost:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"

# Create default buckets
IFS=',' read -ra BUCKETS <<< "$MINIO_DEFAULT_BUCKETS"
for bucket in "${BUCKETS[@]}"; do
  bucket="${bucket#"${bucket%%[![:space:]]*}"}"   # trim leading
  bucket="${bucket%"${bucket##*[![:space:]]}"}"   # trim trailing

  if mc ls "local/$bucket" > /dev/null 2>&1; then
    echo "Bucket '$bucket' already exists."
  else
    echo "Creating bucket: $bucket"
    mc mb "local/$bucket"
  fi
done

# Register Kafka event for PUT events
IFS=',' read -ra KAFKA_BUCKETS <<< "$MINIO_NOTIFY_KAFKA_BUCKETS_1"
for bucket in "${KAFKA_BUCKETS[@]}"; do
  bucket="${bucket#"${bucket%%[![:space:]]*}"}"
  bucket="${bucket%"${bucket##*[![:space:]]}"}"
  
  echo "Removing existing events for bucket: $bucket (if any)"
  mc event remove "local/$bucket" --force || true

  echo "Adding PUT event notification for bucket: $bucket"
  mc event add "local/$bucket" arn:minio:sqs::1:kafka --event put
done

echo "All buckets created and event notifications configured."

# Keep container alive
tail -f /dev/null
