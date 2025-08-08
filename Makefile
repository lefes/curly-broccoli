APP=curly-broccoli

.PHONY: build run tidy fmt vet test docker

build:
	GO111MODULE=on CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/curly-broccoli

run:
	GO111MODULE=on go run ./cmd/curly-broccoli

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet ./...

test:
	go test ./...

docker:
	docker build -t $(APP):latest .

