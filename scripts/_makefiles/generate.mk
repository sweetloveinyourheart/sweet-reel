# Commands for generate everything

goimports: # Fix any errors in the imports with goimports
	@./scripts/lint/go.sh fix-goimports

gen:
	@./scripts/gen/run_code_generators.sh