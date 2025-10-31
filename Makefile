COGNIT_THRESHOLD = 5
CYCLO_THRESHOLD = 5

OUT_PATH:=$(CURDIR)/pkg
LOCAL_BIN:=$(CURDIR)/bin

DATABASE_COMPOSE_FILE = database/db-compose.yaml
KAFKA_COMPOSE_FILE = kafka/kafka-compose.yaml
NOTIFIER_COMPOSE_FILE = docker/notifier/compose.yaml
JAEGER_COMPOSE_FILE = jaeger/jaeger-compose.yaml
E2E_SERVER_PID = .e2e_server.pid
E2E_SERVER_LOG = .e2e_server.log

cognitive-lint:
	@echo "Running cognitive complexity linting..."
	@gocognit -over $(COGNIT_THRESHOLD) internal

cyclomatic-lint:
	@echo "Running cyclomatic complexity linting..."
	@gocyclo -over $(CYCLO_THRESHOLD) internal

lint:
	@make cognitive-lint
	@make cyclomatic-lint
	@echo "Static analysis complete!"

build:
	@echo "Building the application..."
	@go build cmd/hw1/main.go

run:
	@echo "Running the application..."
	@	./main

build-and-run:
	@echo "Building and running the application..."
	@make build
	@make run

run-dev:
	@echo "Running the application..."
	@go run main.go

install-deps:
	@echo "Installing dependencies..."
	@go mod download

coverage:
	@echo "Running test coverage..."
	@go test -coverprofile /dev/null ./internal/usecases ./internal/usecases/packager

run-db:
	@echo "Running the database..."
	@docker compose -f $(DATABASE_COMPOSE_FILE) up -d

compose-up:
	@echo "Running the database and kafka..."
	@docker compose -f $(DATABASE_COMPOSE_FILE) -f $(KAFKA_COMPOSE_FILE) -f $(JAEGER_COMPOSE_FILE) up -d

compose-down:
	@echo "Stopping the database and kafka..."
	@docker compose -f $(DATABASE_COMPOSE_FILE) -f $(KAFKA_COMPOSE_FILE) -f $(JAEGER_COMPOSE_FILE) down

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	@goose -dir ./migrations -s create rename_me sql

goose-up:
	@goose -dir ./migrations postgres ${url} up

goose-down:
	@goose -dir ./migrations postgres ${url} down

squawk-install:
	@npm install -g squawk-cli

squawk-lint:
	@echo "Running squawk linting..."
	@squawk -c .squawk.toml migrations/*.sql

.PHONY: run-prometheus
run-prometheus:
	prometheus --config.file prometheus/prometheus.yaml


# grpc

bin-deps: .vendor-proto
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${OUT_PATH}
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		./api/pvz-service/v1/pvz-service.proto

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp

deps: install-deps bin-deps

all: deps generate build-and-run

# --------------------
# E2E orchestration
# --------------------

.PHONY: e2e-up e2e-test e2e-down e2e

e2e-up: goose-install run-db
	@echo "Applying migrations..."
	@echo "Waiting for Postgres to be healthy..." 
	@for i in $$(seq 1 60); do \
		status=$$(docker inspect --format='{{.State.Health.Status}}' pvz_db 2>/dev/null || echo "starting"); \
		if [ "$$status" = "healthy" ]; then echo "Postgres is healthy"; break; fi; \
		printf "."; sleep 1; \
	done
	@GOOSE_URL='host=localhost port=5430 user=test password=test dbname=test sslmode=disable' && \
		goose -dir ./migrations postgres "$$GOOSE_URL" up
	@echo "Starting API server..."
	@PVZ_ID=PVZ-1 \
		POSTGRES_HOST=localhost \
		POSTGRES_PORT=5430 \
		POSTGRES_USERNAME=test \
		POSTGRES_PASSWORD=test \
		POSTGRES_DATABASE=test \
		nohup go run cmd/grpc-server/main.go > $(E2E_SERVER_LOG) 2>&1 & echo $$! > $(E2E_SERVER_PID)
	@printf "Waiting for HTTP gateway to be ready"
	@for i in $$(seq 1 60); do \
		if curl -sSf http://localhost:8081/metrics >/dev/null 2>&1; then echo "\nGateway is up"; break; fi; \
		printf "."; sleep 0.5; \
		done

e2e-test:
	@echo "Running E2E tests..."
	@E2E=1 go test -v ./test/e2e

e2e-down:
	@echo "Stopping API server and database..."
	@-if [ -f $(E2E_SERVER_PID) ]; then kill $$(cat $(E2E_SERVER_PID)) >/dev/null 2>&1 || true; rm -f $(E2E_SERVER_PID); fi
	@docker compose -f $(DATABASE_COMPOSE_FILE) down

e2e: ## Run full E2E cycle (up -> test -> down)
	@set -e; trap '$(MAKE) e2e-down' EXIT; \
	$(MAKE) e2e-up; \
	$(MAKE) e2e-test