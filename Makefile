.PHONY: test

# Define the name of your binary
BINARY_NAME := url_shortener

# Define flags for the test command
TEST_FLAGS := -covermode=atomic -coverprofile=coverage.out

# Setup docker-compose
docker-compose-up:
	docker-compose up -d

# Default target: run tests
test:
	go test $(TEST_FLAGS) ./...

# Target to run tests with verbose output
test-verbose:
	go test -v ./...

# Target to generate coverage report
coverage:
	go tool cover -html=coverage.out

# Target to clean up generated files
clean:
	go clean
ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
	del /F /Q coverage.out
	del /F /Q $(BINARY_NAME)
else
	rm -f coverage.out
	rm -f $(BINARY_NAME)
endif