#!/bin/bash
. ./scripts/util.sh

set -e

KAFKA_CONTAINER=${KAFKA_CONTAINER:-"srl-kafka"}
KAFKA_BROKER=${KAFKA_BROKER:-"localhost:9092"}

app_echo "Waiting for Kafka to be ready..."
until docker exec "$KAFKA_CONTAINER" /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server "$KAFKA_BROKER" > /dev/null 2>&1; do
  app_echo "Kafka is not ready yet. Retrying in 5 seconds..."
  sleep 5
done

app_echo "Kafka is ready. Creating topics..."

# Function to create topic if it doesn't exist
create_topic() {
    local topic_name=$1
    local partitions=${2:-3}
    local replication=${3:-1}

    app_echo "Creating topic: $topic_name (partitions=$partitions, replication=$replication)"

    docker exec "$KAFKA_CONTAINER" /opt/kafka/bin/kafka-topics.sh --bootstrap-server "$KAFKA_BROKER" \
        --create --if-not-exists \
        --topic "$topic_name" \
        --partitions "$partitions" \
        --replication-factor "$replication" \
        --config retention.ms=604800000 \
        --config segment.ms=86400000 || app_echo "Topic '$topic_name' may already exist"
}

# CDC topics from configuration file
CDC_TOPICS_CONFIG="./scripts/cdc/kafka/cdc-topics.conf"

if [ -f "$CDC_TOPICS_CONFIG" ]; then
    app_echo "Reading CDC topics from $CDC_TOPICS_CONFIG"
    while IFS=':' read -r topic_name partitions replication || [ -n "$topic_name" ]; do
        # Skip comments and empty lines
        [[ "$topic_name" =~ ^#.*$ ]] || [ -z "$topic_name" ] && continue

        # Use defaults if not specified
        partitions=${partitions:-3}
        replication=${replication:-1}

        create_topic "$topic_name" "$partitions" "$replication"
    done < "$CDC_TOPICS_CONFIG"
else
    app_echo "Warning: CDC topics configuration file not found at $CDC_TOPICS_CONFIG"
fi

app_echo ""
app_echo "Listing all topics:"
docker exec "$KAFKA_CONTAINER" /opt/kafka/bin/kafka-topics.sh --bootstrap-server "$KAFKA_BROKER" --list

app_echo ""
app_echo "AAll topics created successfully!"
