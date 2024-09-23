COGNIT_THRESHOLD = 5
CYCLO_THRESHOLD = 5

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
	@make lint
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
