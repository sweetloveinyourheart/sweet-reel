# Commands for setup cdc with debezium and elasticsearch

setup-es-indices:
	@./scripts/cdc/init_indices.sh

setup-es-connectors:
	@./scripts/cdc/init_connectors.sh

cdc-setup:
	@make setup-es-indices
	@make setup-es-connectors