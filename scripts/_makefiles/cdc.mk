# Commands for setup cdc with debezium and elasticsearch

setup-dbz-indices:
	@./scripts/cdc/init_indices.sh

setup-es-connectors:
	@./scripts/cdc/init_connectors.sh

cdc-setup:
	@make setup-dbz-indices
	@make setup-es-connectors