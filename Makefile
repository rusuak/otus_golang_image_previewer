ifeq ($(OS),Windows_NT)
    GOOS := windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        GOOS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        GOOS := darwin
    endif
endif

.PHONY: build
build:
	GOOS=${GOOS} CGO_ENABLED=0 GOARCH=amd64 go build -v -trimpath -o ./server ./cmd/main.go

.PHONY: run
run: docker-build docker-up

.PHONY: test
test:
	go test ./internal/... ./pkg/... -v -race -count=10

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker-build
docker-build:
	docker-compose -f docker-compose.yaml build

.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up -d --remove-orphans

.PHONY: docker-down
docker-down:
	docker-compose -f docker-compose.yaml down --remove-orphans
