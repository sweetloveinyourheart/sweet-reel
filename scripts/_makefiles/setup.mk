# Targets for setting up the local development environment

go-deps: # Install dependencies for Go
	@./scripts/deps/install_go_deps.sh || ./scripts/deps/install_go_deps.sh

# Commands for setup cdc with debezium and elasticsearch in local development
cdc-setup:
	@./scripts/development/init_topics.sh
	@./scripts/development/init_indices.sh
	@./scripts/development/init_connectors.sh
