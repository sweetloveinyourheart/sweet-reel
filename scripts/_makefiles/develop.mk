# Commands for local development
# Targets for local development and testing
FULL_SERVER_STACK_COMPOSE_FILE := ./dockerfiles/docker-compose.yml

base-compose-up:
	@source ./scripts/util.sh && app_compose_up "$(COMPOSE_FILE)"

base-compose-down:
	@source ./scripts/util.sh && app_compose_down "$(COMPOSE_FILE)"

compose-up: # Start the full-server stack
	@make base-compose-up COMPOSE_FILE=$(FULL_SERVER_STACK_COMPOSE_FILE)

compose-down: 
	@make base-compose-down COMPOSE_FILE=$(FULL_SERVER_STACK_COMPOSE_FILE)
