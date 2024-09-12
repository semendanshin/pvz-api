
cognitive-lint:
	@echo "Running cognitive complexity linting..."
	@gocognit -over 5 .

cyclomatic-lint:
	@echo "Running cyclomatic complexity linting..."
	@gocyclo -over 5 .

lint: cognitive-lint cyclomatic-lint

build:
	@echo "Building the application..."
	lint
	@go build -o bin/ ./...
