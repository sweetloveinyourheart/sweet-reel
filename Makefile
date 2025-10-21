SHELL := /bin/bash

include $(PWD)/scripts/_makefiles/build.mk
include $(PWD)/scripts/_makefiles/develop.mk
include $(PWD)/scripts/_makefiles/generate.mk
include $(PWD)/scripts/_makefiles/lint.mk
include $(PWD)/scripts/_makefiles/setup.mk
include $(PWD)/scripts/_makefiles/tests_unit.mk

help: # Print this help message
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

env: # Print makefile's environment (different from your local shell)
	@printenv
