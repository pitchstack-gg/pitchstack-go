.DEFAULT_GOAL := ci

.PHONY: format lint gosec build test-unit test-integration test ci

format:
	@echo "==> formatting"
	go fmt ./...

lint:
	@echo "==> linting"
	golangci-lint run ./...

gosec:
	@echo "==> running gosec"
	gosec ./...

test-unit:
	go test ./client/...

test-integration:
	@echo "==> running integration tests"
	@if [ -z "$$PITCHSTACK_API_TOKEN" ]; then \
		echo "PITCHSTACK_API_TOKEN not set; skipping integration tests"; \
	else \
		go test -tags=integration ./tests/integration; \
	fi

test: test-unit

ci: format lint gosec test-unit test-integration
