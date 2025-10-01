# Targets for setting up the local development environment

go-deps: # Install dependencies for Go
	@./scripts/deps/install_go_deps.sh || ./scripts/deps/install_go_deps.sh