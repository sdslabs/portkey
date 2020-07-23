# Prints help message
help:
	@echo "Portkey"
	@echo "build   - Build portkey"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"

# Build portkey
build:
	@./scripts/build/build.sh

# Format code using golangci-lint
format:
	@./scripts/build/format.sh

# Install required tools
install:
	@./scripts/build/install.sh

# Lint code using golangci-lint
lint:
	@./scripts/build/lint.sh 

.PHONY: help build format install lint