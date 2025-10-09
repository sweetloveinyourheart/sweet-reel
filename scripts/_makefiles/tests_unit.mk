# Targets for running unit and integration tests under `go test`

test: # Run all unit tests (see more options in Makefile)
	@./scripts/unit_test/run_all_unit_tests.sh
test-verbose:
	@./scripts/unit_test/run_all_unit_tests.sh verbose
test-coverage:
	@./scripts/unit_test/run_all_unit_tests.sh cov
	@./scripts/unit_test/print_coverage_stats.sh
test-coverage-verbose:
	@./scripts/unit_test/run_all_unit_tests.sh cov verbose
	@./scripts/unit_test/print_coverage_stats.sh
print-coverage:
	@./scripts/unit_test/print_coverage_stats.sh


# CI Automation Conventions:
# Any makefile target that starts with ut- will run the unit tests for that package
# Any makefile target that starts with cov- will run the unit tests for that package and generate a coverage report
#
# Any unit tests that are not covered by an explicit ut-/cov- target will be covered by 'ut-other' and 'cov-other.'
# If you create a new pair of ut- and cov- targets, remember to exclude that package from 'template-other', and add a new job to .github/workflows/tests.yaml

template-ut:
	@go clean -testcache
	@./scripts/unit_test/ci_test_wrapper.sh "$(optionalArg)" "$(package)" "$(verbose)" || exit 1

template-cov:
	@rm -rf tests/logs/cov-$(packageName)*
	@mkdir -p tests/logs/cov-$(packageName)
	@(make template-ut package=$(package) packageName=$(packageName) verbose=$(verbose) optionalArg="-coverprofile=tests/logs/cov-$(packageName)/cov.tmp") || exit 1
	@exclusions=$$(grep --include=\*.go -Ril "DO NOT EDIT" . | cut -c 3- | xargs | tr -s '[:blank:]' ',' | sed -E 's!,!|github.com/sweetloveinyourheart/exploding-kittens/!g'); \
	cat tests/logs/cov-$(packageName)/cov.tmp | grep -vE "github.com/sweetloveinyourheart/exploding-kittens/$${exclusions}" > tests/logs/cov-$(packageName)/cov;
	@rm -f tests/logs/cov-$(packageName)/cov.tmp
	@go tool cover -func tests/logs/cov-$(packageName)/cov                                                  >> tests/logs/cov-$(packageName)/low-level.txt  || exit 1
	@go tool cover -func tests/logs/cov-$(packageName)/cov | grep total: | awk '{print $$3}' | sed 's/.$$//' > tests/logs/cov-$(packageName)/percentage.txt || exit 1
	@go tool cover -html=tests/logs/cov-$(packageName)/cov -o                                                  tests/logs/cov-$(packageName)/visual.html    || exit 1
	@echo Reports are available in the logs directory.

###
### SERVICES
###

ut-video_management:
	@make template-ut package=services/video_management packageName=video_management

cov-video_management:
	@make template-cov package=services/video_management packageName=video_management

ut-video_processing:
	@make template-ut package=services/video_processing packageName=video_processing

cov-video_processing:
	@make template-cov package=services/video_processing packageName=video_processing

ut-user:
	@make template-ut package=services/user packageName=user

cov-user:
	@make template-cov package=services/user packageName=user