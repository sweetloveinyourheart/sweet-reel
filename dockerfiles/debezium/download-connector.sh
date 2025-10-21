#!/bin/bash
set -e

CONNECTOR_DIR="/kafka/connect/confluentinc-kafka-connect-elasticsearch"
mkdir -p "$CONNECTOR_DIR"
cd "$CONNECTOR_DIR"

# Download the connector package from a mirror or build from source
# Since direct package downloads are problematic, we'll manually download all required JARs

echo "Downloading Elasticsearch connector JARs..."

# Main connector JAR from Maven Central (this might be a thin JAR, so we need all dependencies)
curl -f -o kafka-connect-elasticsearch-14.0.10.jar \
  "https://repo1.maven.org/maven2/io/confluent/kafka-connect-elasticsearch/14.0.10/kafka-connect-elasticsearch-14.0.10.jar"

# Check if it's a valid JAR (more than 10KB)
if [ $(stat -f%z kafka-connect-elasticsearch-14.0.10.jar 2>/dev/null || stat -c%s kafka-connect-elasticsearch-14.0.10.jar) -lt 10000 ]; then
    echo "ERROR: Downloaded JAR is too small, likely not the actual connector"
    rm -f kafka-connect-elasticsearch-14.0.10.jar
    exit 1
fi

echo "Connector downloaded successfully!"
