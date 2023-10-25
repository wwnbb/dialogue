# Makefile with Colored Help and Virtual Environment Activation

# Set the default target as 'help'
.DEFAULT_GOAL := help

# Define the colors to be used in the help text
ifndef fish
    YELLOW := $(shell tput -Txterm setaf 3)
    GREEN := $(shell tput -Txterm setaf 2)
    RESET := $(shell tput -Txterm sgr0)
else
    set -x YELLOW (set_color yellow)
    set -x GREEN (set_color green)
    set -x RESET (set_color normal)
endif

# Define the targets and their commands
.PHONY: help
help:
	@echo "$(YELLOW)Usage:$(RESET)"
	@echo "  make $(GREEN)<target>$(RESET)"
	@echo ""
	@echo "$(YELLOW)Targets:$(RESET)"
	@echo "  $(GREEN)build$(RESET)          Build gophish application"
	@echo "  $(GREEN)debug$(RESET)          Build debug version of gophish application"
	@echo "  $(GREEN)test$(RESET)           Run tests"
	@echo "  $(GREEN)tidy$(RESET)           Run go mod tidy"

.PHONY: build
build:
	@go build .

.PHONY: debug
debug:
	@go build -gcflags=all="-N -l" .

.PHONY: test
test:
	@go test ./...

.PHONY: tidy
tidy:
	@go mod tidy
