# CDC and Elasticsearch Setup Guide

This document describes the Change Data Capture (CDC) setup using Debezium and Elasticsearch integration for the Sweet Reel project.

## Architecture Overview

The CDC pipeline captures database changes from PostgreSQL and streams them to Elasticsearch for real-time search capabilities:

```
PostgreSQL (WAL) → Debezium Connect → Kafka → Elasticsearch
```

### Components

1. **PostgreSQL with Logical Replication**: Configured with `wal_level=logical` to enable CDC
2. **Debezium Connect**: Captures database changes and publishes to Kafka topics
3. **Apache Kafka**: Message broker for change events
4. **Elasticsearch**: Search engine for indexed data

## Services

### Elasticsearch
- **Host**: `elasticsearch` (internal), `localhost:9200` (external)
- **Ports**:
  - 9200 (REST API)
  - 9300 (Transport)
- **IP**: 172.16.244.71
- **Data Volume**: `elasticsearch_data`

### Debezium Connect
- **Host**: `debezium` (internal), `localhost:8083` (external)
- **Port**: 8083 (REST API)
- **IP**: 172.16.244.81
- **Configuration**: Stored in Kafka topics

## Database Configuration

PostgreSQL is configured with the following parameters for CDC:

- `wal_level=logical` - Enable logical replication
- `max_wal_senders=10` - Maximum concurrent WAL senders
- `max_replication_slots=10` - Maximum replication slots

## Indices

Elasticsearch indices are created automatically from the corresponding PostgreSQL tables captured by Debezium connectors.

**Configuration**: See `scripts/cdc/elasticsearch/indices.conf` for the list of indices and their mapping files.

**Security Note**: Sensitive tables containing authentication credentials (OAuth tokens, passwords, etc.) should be excluded from CDC.

## Debezium Connectors

Debezium source connectors capture changes from PostgreSQL tables and publish them to Kafka topics.

**Configuration**: See `scripts/cdc/debezium/connectors.conf` for the list of connectors to register.

**Topic Naming**: All CDC topics use the `cdc-` prefix with hyphens (following Kafka best practices) to differentiate them from application topics.

## Setup Instructions

### 1. Start the Services

```bash
# Start all services including Elasticsearch and Debezium
make compose-up
```

This will:
- Start PostgreSQL with logical replication enabled
- Start Elasticsearch
- Start Debezium Connect
- Wait for all services to be healthy

### 2. Initialize Kafka Topics

Create the required Kafka topics before registering connectors:

```bash
# From the project root
make setup-dbz-topics
```

This creates all CDC topics defined in `scripts/cdc/kafka/cdc-topics.conf` (prefixed with `cdc-`) plus Debezium Connect internal topics.

### 3. Initialize Elasticsearch Indices

The Elasticsearch indices need to be created before data can be synced. Run the initialization script:

```bash
# From the project root
make setup-dbz-indices
```

### 4. Register Debezium Connectors

Register the CDC connectors to start capturing database changes:

```bash
# From the project root
make setup-es-connectors
```

### 5. Verify Setup

#### Check Elasticsearch Health
```bash
curl http://localhost:9200/_cluster/health?pretty
```

#### List Elasticsearch Indices
```bash
curl http://localhost:9200/_cat/indices?v
```

#### Check Debezium Connectors
```bash
curl http://localhost:8083/connectors
```

#### Check Connector Status
```bash
# List all connectors
curl http://localhost:8083/connectors | jq '.'

# Check specific connector status (replace <connector-name>)
curl http://localhost:8083/connectors/<connector-name>/status | jq '.'
```

#### List Kafka Topics
```bash
docker exec srl-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Usage Examples

### Search Documents

```bash
curl -X GET "http://localhost:9200/<index-name>/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "match": {
        "<field-name>": "<search-term>"
      }
    }
  }'
```

### Get Document by ID

```bash
curl -X GET "http://localhost:9200/<index-name>/_doc/<document-id>?pretty"
```

## Monitoring

### View CDC Events in Kafka

```bash
# Monitor a specific CDC topic (replace <topic-name>)
docker exec srl-kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic <topic-name> \
  --from-beginning
```

### Check Elasticsearch Document Count

```bash
# Check document count for an index (replace <index-name>)
curl -X GET "http://localhost:9200/<index-name>/_count?pretty"
```

### View Debezium Connect Logs

```bash
docker logs -f srl_debezium
```

## Troubleshooting

### Unknown Topic or Partition Warning

If you see warnings like:
```
WARN Error while fetching metadata with correlation id X : {cdc-<table-name>=UNKNOWN_TOPIC_OR_PARTITION}
```

This means the Kafka topics haven't been created yet. Run:
```bash
make setup-dbz-topics
```

Then restart the affected connector:
```bash
curl -X POST http://localhost:8083/connectors/<connector-name>/restart
```

### Connectors Not Starting

1. Check Debezium logs:
   ```bash
   docker logs srl_debezium
   ```

2. Verify database permissions:
   ```bash
   docker exec srl_database psql -U root_admin -d <database-name> -c "SELECT * FROM pg_replication_slots;"
   ```

3. Restart the connector:
   ```bash
   curl -X POST http://localhost:8083/connectors/<connector-name>/restart
   ```

### Missing Data in Elasticsearch

1. Check if connector is running:
   ```bash
   curl http://localhost:8083/connectors/<connector-name>/status
   ```

2. Verify Kafka topics have data:
   ```bash
   docker exec srl-kafka /opt/kafka/bin/kafka-console-consumer.sh \
     --bootstrap-server localhost:9092 \
     --topic <topic-name> \
     --from-beginning \
     --max-messages 10
   ```

3. Check Elasticsearch indices exist:
   ```bash
   curl http://localhost:9200/_cat/indices?v
   ```

### Replication Slot Issues

If you need to reset the replication slot:

```bash
# Connect to PostgreSQL
docker exec -it srl_database psql -U root_admin -d <database-name>

# Drop the replication slot (check connector config for slot name)
SELECT pg_drop_replication_slot('<slot-name>');

# Exit psql
\q
```

Then delete and recreate the Debezium connector.

### Delete and Recreate Connector

```bash
# Delete connector
curl -X DELETE http://localhost:8083/connectors/<connector-name>

# Wait a few seconds, then recreate
curl -X POST http://localhost:8083/connectors \
  -H 'Content-Type: application/json' \
  -d @scripts/cdc/debezium/<connector-name>.json
```

## Data Consistency

### Initial Snapshot

When a Debezium connector starts for the first time, it:
1. Takes a consistent snapshot of all existing data
2. Streams the snapshot to Kafka topics
3. Switches to streaming mode for ongoing changes

### Change Types

Debezium captures the following change types:
- **CREATE** (c) - New row inserted
- **UPDATE** (u) - Existing row updated
- **DELETE** (d) - Row deleted

### Handling Deletes

When a row is deleted in PostgreSQL, Debezium sends a tombstone event to Kafka. The Elasticsearch sink connector is configured to delete the corresponding document.

## Performance Considerations

1. **Elasticsearch Memory**: Default is 512MB heap. Increase if needed in docker-compose.yml:
   ```yaml
   ES_JAVA_OPTS: "-Xms1g -Xmx1g"
   ```

2. **Kafka Retention**: Change events are stored in Kafka topics. Configure retention as needed.

3. **Replication Lag**: Monitor the lag between PostgreSQL changes and Elasticsearch updates.

4. **Index Refresh**: Elasticsearch indices refresh every second by default. Adjust for your use case.

## Security Considerations

Current setup is for development only. For production:

1. Enable Elasticsearch security (xpack.security)
2. Use SSL/TLS for all connections
3. Secure Debezium Connect REST API
4. Use secrets management for passwords
5. Restrict network access
6. Enable authentication on Kafka

## Further Reading

- [Debezium PostgreSQL Connector](https://debezium.io/documentation/reference/stable/connectors/postgresql.html)
- [Elasticsearch Guide](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [PostgreSQL Logical Replication](https://www.postgresql.org/docs/current/logical-replication.html)
- [Kafka Connect](https://kafka.apache.org/documentation/#connect)
