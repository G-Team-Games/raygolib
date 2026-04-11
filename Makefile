.PHONY: help test coverage clean-coverage sonar-up sonar-down sonar-logs sonar-scan

DOCKER_COMPOSE ?= docker compose
COMPOSE_PATH = dev/docker-compose.yml
COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_FILTERED := $(COVERAGE_DIR)/filtered.out
COVERAGE_EXCL := internal/raylib internal/testutils cmd/covfilter

help:
	@echo "Available targets:"
	@echo "  make test            - run all Go tests"
	@echo "  make coverage        - generate coverage report"
	@echo "  make clean-coverage  - remove coverage artifacts"
	@echo "  make sonar-up        - start SonarQube"
	@echo "  make sonar-down      - stop SonarQube"
	@echo "  make sonar-logs      - tail SonarQube logs"
	@echo "  make sonar-scan      - run sonar scanner (requires SONAR_TOKEN in environment)"

test:
	go test -race ./...

coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_OUT) ./...
	go run ./cmd/covfilter $(COVERAGE_EXCL) < $(COVERAGE_OUT) > $(COVERAGE_FILTERED)
	go tool cover -func=$(COVERAGE_FILTERED)

clean-coverage:
	rm -rf $(COVERAGE_DIR)

sonar-up:
	$(DOCKER_COMPOSE) -f $(COMPOSE_PATH) up -d sonarqube

sonar-down:
	$(DOCKER_COMPOSE) down

sonar-logs:
	$(DOCKER_COMPOSE) logs -f sonarqube

sonar-scan: coverage
	$(DOCKER_COMPOSE) -f $(COMPOSE_PATH) run --rm sonar-scanner
