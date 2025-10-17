# Commands for setup cdc with debezium and elasticsearch

setup-dbz-topics:
	@./scripts/cdc/init_topics.sh

setup-dbz-indices:
	@./scripts/cdc/init_indices.sh

setup-es-connectors:
	@./scripts/cdc/init_connectors.sh

cdc-setup:
	@make setup-dbz-topics
	@make setup-dbz-indices
	@make setup-es-connectors