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

The following Elasticsearch indices are created automatically:

### Video Management Indices

1. **videos** - Main video metadata
   - Fields: id, uploader_id, title, description, status, object_key, processed_at, created_at, updated_at
   - Full-text search on: title, description

2. **video_manifests** - HLS manifest files
   - Fields: id, video_id, object_key, size_bytes, created_at

3. **video_variants** - Video quality variants
   - Fields: id, video_id, quality, object_key, total_segments, total_duration, created_at

4. **video_thumbnails** - Video thumbnail metadata
   - Fields: id, video_id, object_key, width, height, created_at

### User Indices

5. **users** - User profiles
   - Fields: id, email, name, picture, created_at, updated_at
   - Full-text search on: name

6. **user_identities** - OAuth provider identities
   - Fields: id, user_id, provider, provider_user_id, access_token (not indexed), refresh_token (not indexed), expires_at, created_at, updated_at

## Debezium Connectors

### Video Management Connector
- **Name**: `video-management-connector`
- **Database**: `video-management`
- **Tables**: videos, video_manifests, video_variants, video_thumbnails
- **Slot**: `debezium_video_management`
- **Topics**: videos, video_manifests, video_variants, video_thumbnails

### User Connector
- **Name**: `user-connector`
- **Database**: `user`
- **Tables**: users, user_identities
- **Slot**: `debezium_user`
- **Topics**: users, user_identities

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

### 2. Initialize Elasticsearch Indices

The Elasticsearch indices need to be created before data can be synced. Run the initialization script:

```bash
# From the project root
make setup-es-indices
```

### 3. Register Debezium Connectors

Register the CDC connectors to start capturing database changes:

```bash
# From the project root
make setup-es-connectors
```

### 4. Verify Setup

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
curl http://localhost:8083/connectors/video-management-connector/status | jq '.'
curl http://localhost:8083/connectors/user-connector/status | jq '.'
```

#### List Kafka Topics
```bash
docker exec srl-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Usage Examples

### Search Videos by Title

```bash
curl -X GET "http://localhost:9200/videos/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "match": {
        "title": "tutorial"
      }
    }
  }'
```

### Search Users by Name

```bash
curl -X GET "http://localhost:9200/users/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "match": {
        "name": "john"
      }
    }
  }'
```

### Get Video by ID

```bash
curl -X GET "http://localhost:9200/videos/_doc/{VIDEO_UUID}?pretty"
```

### Get All Videos by Uploader

```bash
curl -X GET "http://localhost:9200/videos/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "term": {
        "uploader_id": "USER_UUID"
      }
    }
  }'
```

### Filter Videos by Status

```bash
curl -X GET "http://localhost:9200/videos/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "term": {
        "status": "ready"
      }
    }
  }'
```

## Monitoring

### View CDC Events in Kafka

```bash
# Monitor videos topic
docker exec srl-kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic videos \
  --from-beginning
```

### Check Elasticsearch Document Count

```bash
curl -X GET "http://localhost:9200/videos/_count?pretty"
curl -X GET "http://localhost:9200/users/_count?pretty"
```

### View Debezium Connect Logs

```bash
docker logs -f srl_debezium
```

## Troubleshooting

### Connectors Not Starting

1. Check Debezium logs:
   ```bash
   docker logs srl_debezium
   ```

2. Verify database permissions:
   ```bash
   docker exec srl_database psql -U root_admin -d video-management -c "SELECT * FROM pg_replication_slots;"
   ```

3. Restart the connector:
   ```bash
   curl -X POST http://localhost:8083/connectors/video-management-connector/restart
   ```

### Missing Data in Elasticsearch

1. Check if connector is running:
   ```bash
   curl http://localhost:8083/connectors/video-management-connector/status
   ```

2. Verify Kafka topics have data:
   ```bash
   docker exec srl-kafka /opt/kafka/bin/kafka-console-consumer.sh \
     --bootstrap-server localhost:9092 \
     --topic videos \
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
docker exec -it srl_database psql -U root_admin -d video-management

# Drop the replication slot
SELECT pg_drop_replication_slot('debezium_video_management');

# Exit psql
\q
```

Then restart the Debezium connector.

### Delete and Recreate Connector

```bash
# Delete connector
curl -X DELETE http://localhost:8083/connectors/video-management-connector

# Wait a few seconds, then recreate
curl -X POST http://localhost:8083/connectors \
  -H 'Content-Type: application/json' \
  -d @dockerfiles/debezium/video-management-connector.json
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
