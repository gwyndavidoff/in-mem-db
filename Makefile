.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run src/main.go