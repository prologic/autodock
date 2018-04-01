.PHONY: dev build image test deps clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
APP=autodock
REPO?=prologic/$(APP)
TAG?=latest
BUILD?=-dev

all: dev

dev: build
	@./cmd/$(APP)/$(APP) -debug

deps:
	@go get ./...

build: clean
	@echo " -> Building $(TAG)$(BUILD)"
	@cd cmd/$(APP) && go build -tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT) -X github.com/$(REPO)/version.Build=$(BUILD)" .
	@echo "Built $$(./cmd/$(APP)/$(APP) -v)"

image:
	@docker build --build-arg TAG=$(TAG) --build-arg BUILD=$(BUILD) -t $(REPO):$(TAG) .
	@echo "Image created: $(REPO):$(TAG)"

test:
	@go test -v -cover -race $(TEST_ARGS) ./...

clean:
	@rm -rf $(APP)
