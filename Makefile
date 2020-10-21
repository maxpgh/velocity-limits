.PHONY: test
test:
	go test -v ./... -cover -count=1

.PHONY: build
build:
	go build -o ./bin/velocity-limits -race -ldflags '-w -s' .

.PHONY: run
run:
	go run main.go