FROM debezium/connect:2.5

USER root

# Download and extract the full Elasticsearch connector from Confluent archive
RUN cd /tmp && \
    curl -L -o es-connector.tar.gz "https://github.com/jcustenborder/kafka-connect-archive/raw/master/confluentinc-kafka-connect-elasticsearch-11.1.3.tar.gz" && \
    mkdir -p /kafka/connect && \
    tar -xzf es-connector.tar.gz -C /kafka/connect/ && \
    rm es-connector.tar.gz && \
    chmod -R 755 /kafka/connect/confluentinc-kafka-connect-elasticsearch

USER kafka
