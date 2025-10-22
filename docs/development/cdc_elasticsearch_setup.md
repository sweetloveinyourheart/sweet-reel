# CDC and Elasticsearch Setup Guide

Change Data Capture (CDC) setup using Debezium and Elasticsearch for real-time search capabilities.

## Architecture

```
PostgreSQL → Debezium → Kafka → Elasticsearch Sink → Elasticsearch
              (unwrap)    ↓         (extract ID)         ↓
                      Clean data                   Clean documents
```

**Components:**
- PostgreSQL (logical replication enabled)
- Debezium Connect (CDC capture + transforms)
- Kafka (message broker)
- Confluent Elasticsearch Sink Connector
- Elasticsearch (search engine)

## Setup

### 1. Install Confluent Elasticsearch Connector

**Required:** Download manually from [Confluent Hub](https://www.confluent.io/hub/confluentinc/kafka-connect-elasticsearch)

Ensure docker-compose.yml mounts the plugin directory:
```yaml
# debezium:
#   volumes:
#     - ./debezium/connectors:/kafka/connect/debezium-connector-elasticsearch
```

### 2. Start Services

```bash
make compose-up
```

### 3. Initialize CDC

```bash
make cdc-setup
```

This runs:
1. Creates Elasticsearch indices with mappings
2. Registers Debezium connectors (topics auto-created)

## Key Points

- **Topics auto-created** by Debezium source connectors
- **Clean data** via `ExtractNewRecordState` transform on source side
- **Document ID** extracted from `id` field via `ValueToKey` + `ExtractField$Key` transforms
- **Delete handling:** Tombstone events delete documents in Elasticsearch
- **Development only:** No security enabled (enable for production)

## Resources

- [Debezium PostgreSQL Connector](https://debezium.io/documentation/reference/stable/connectors/postgresql.html)
- [Confluent Elasticsearch Sink](https://docs.confluent.io/kafka-connectors/elasticsearch/current/overview.html)
- [Kafka Connect Transformations](https://docs.confluent.io/platform/current/connect/transforms/overview.html)
