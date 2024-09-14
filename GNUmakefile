default: test

.PHONY: testall test lint docs

test:
	go test ./...

acceptance:
	TEST_ACC=1 go test ./...

lint:
	golangci-lint run ./...
