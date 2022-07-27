# Build system related variables
GOCMD = go
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

# Project related variables
WORKING_DIR = $(CURDIR)
BUILD_DIR = $(WORKING_DIR)/_bin

# build related variables
_SEPARATOR = _
_EXE_POSTFIX = 

ifeq ($(GOOS),windows)
	_EXE_POSTFIX = .exe
endif

# some globally assembled variables
APPLICATION_NAME = sclix_woof
PLATFORM_STRING = $(GOOS)$(_SEPARATOR)$(GOARCH)
EXECUTABLE_NAME = $(APPLICATION_NAME)$(_SEPARATOR)$(PLATFORM_STRING)$(_EXE_POSTFIX)

# some make file variables
LOG_PREFIX = --

.PHONY: clean
clean:
	@$(GOCMD) clean
	@rm -f -r $(BUILD_DIR)

.PHONY: format
format:
	gofmt -w -l -e .

.PHONY: lint
lint:
	./scripts/lint.sh

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/$(EXECUTABLE_NAME): $(BUILD_DIR)
	@echo "Building..."
	@echo "$(LOG_PREFIX) Building ( $(BUILD_DIR)/$(EXECUTABLE_NAME) )"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build -o $(BUILD_DIR)/$(EXECUTABLE_NAME) $(WORKING_DIR)/cmd/main.go

.PHONY: build
build: $(BUILD_DIR)/$(EXECUTABLE_NAME)

.PHONY: test
test: 
	@echo "Testing..."

.PHONY: help
help:
	@echo "Main targets:"
	@echo "$(LOG_PREFIX) format"
	@echo "$(LOG_PREFIX) lint"
	@echo "$(LOG_PREFIX) build"
	@echo "$(LOG_PREFIX) test"
	@echo "$(LOG_PREFIX) clean"
	@echo "\nAvailable parameter:"
	@echo "$(LOG_PREFIX) GOOS                       Specify Operating System to compile for (see golang GOOS, default=$(GOOS))"
	@echo "$(LOG_PREFIX) GOARCH                     Specify Architecture to compile for (see golang GOARCH, default=$(GOARCH))"
