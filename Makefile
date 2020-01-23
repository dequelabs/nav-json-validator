SRC := main.go $(wildcard navjson/*.go)
TESTS := $(wildcard navjson/*_test.go)

nav-json-validator: $(SRC)
	go build

clean:
	rm -f nav-json-validator coverage.out
.PHONY: clean

test:
	go test ./... -v -race
.PHONY: test

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

coverage: coverage.out
	go tool cover -html=$<
.PHONY: coverage

coverage.out: $(SRC) $(TESTS)
	go test ./... -coverprofile=$@