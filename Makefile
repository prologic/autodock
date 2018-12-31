.PHONY: dev build image profile bench test clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
APP=autodock
REPO?=prologic/$(APP)
TAG?=latest
BUILD?=-dev

all: dev

dev: build
	@./cmd/$(APP)/$(APP) -debug

build: clean
	@echo " -> Building $(TAG)$(BUILD)"
	@cd cmd/$(APP) && go build -tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT) -X github.com/$(REPO)/version.Build=$(BUILD)" .
	@echo "Built $$(./cmd/$(APP)/$(APP) -v)"

image:
	@docker build --build-arg TAG=$(TAG) --build-arg BUILD=$(BUILD) -t $(REPO):$(TAG) .
	@echo "Image created: $(REPO):$(TAG)"

profile:
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench ./...

bench:
	@go test -v -bench ./...

test:
	@go test -v -cover -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... -race ./...

clean:
	@git clean -f -d -X
