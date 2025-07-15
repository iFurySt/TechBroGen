.PHONY: run build clean test dev

# Default target
dev: run

# Run the application in development mode
run:
	go run main.go

# Build the application
build:
	go build -o tbg main.go

# Clean build artifacts
clean:
	rm -f tbg

# Run tests (if any)
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy

# Run production build
prod: build
	./tbg