.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. ./...

.PHONY: run
run:
	go run src/main.go