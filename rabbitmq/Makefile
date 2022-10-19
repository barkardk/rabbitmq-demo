.DEFAULT_GOAL := help
BUILD_TIME := $(shell date -u '+%F_%T')
VERSION ?= $(shell (git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0) | cut -c2-)
TAG ?= ${VERSION}
DOCKER_REGISTRY?=ghcr.io/barkardk
.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z\\._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

docker.restart: docker.stop docker.start ## restart docker-compose
build: fmt lint dep imp build.linux ## Run a full build with linting , depencency testing and import
release: docker-manifest.purge publish pi manifests ## Build linux release containers
pies: docker-manifest.purge pi manifests
release.binary: build.linux  ## Build binaries for various linux architectures

docker.start: ## Start docker-compose
	docker-compose up -d --remove-orphans;
docker.stop: ## Stop docker-compose
	docker-compose stop;
docker.build: ## Build rabbitmq-client docker container
	docker build -t rabbitmq-client

.PHONY: clean
clean: ## Clean compiled binaries
	@echo "-> $@"
	@rm -rf target

# -- build local multiarch binaries --
TARGETARCH := amd64 arm arm64
.PHONY: build.linux
build.linux: $(TARGETARCH) ## Build multiarch linux binaries

.PHONY: $(TARGETARCH)
ARCH = ${word 1, $@}
$(TARGETARCH): ## Build multiarch linux binaries
	@echo "-> $@"
	GO111MODULE=on CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux go build -o target/linux-$(ARCH)/mq_test -tags ${VERSION} -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" ./it
	@chmod 755 target/linux-$(ARCH)/mq_test

docker-manifest.purge:
	@docker manifest push --purge ${DOCKER_REGISTRY}/rabbitmq-client:latest || true
# -- build and push older arm containers --
TARGETS ?= linux/arm/v7 linux/arm/v6 linux/arm64/v8
temp = $(subst /, ,$@)
TARGETPLATFORM=$(word 1, $@)
GOOS= $(word 1, $(temp))
GOARCH=$(word 1, $(temp))
ARMV=$(word 3, $(temp))
.PHONY: $(TARGETS) pi
pi: $(TARGETS)
$(TARGETS):
	@echo "-> $@"
	@echo "-- Build containers for VERSION:${VERSION} TAG:${TAG} --"
	@echo "1 $(GOOS) 2 $(GOARCH) 3 -$(ARMV)"
	docker  buildx build --platform $(TARGETPLATFORM) -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-${TAG} -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-latest --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg VERSION=${VERSION} --push .
	@echo "-- Generate manifests for multiarch support --"
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:${TAG}" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-${TAG}"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:${TAG}"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-${TAG}" --os=$(GOOS) --arch=$(GOARCH) --variant=$(ARMV)
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:latest" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-latest"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:latest"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-$(ARMV)-latest" --os=$(GOOS) --arch=$(GOARCH) --variant=$(ARMV)

# -- build and push docker containers --
PLATFORMS ?= linux/amd64 linux/386 linux/arm linux/arm64
TPLATFORM=$(word 1, $@)
temp = $(subst /, ,$@)
GOOS = $(word 1, $(temp))
GOARCH = $(word 2, $(temp))

.PHONY: $(PLATFORMS) publish
publish: $(PLATFORMS) ## Build docker multiarch containers and push to registry
$(PLATFORMS):
	@echo "-> $@"
	@echo "-- Build containers for VERSION:${VERSION} TAG:${TAG} --"
	docker  buildx build --platform $(TPLATFORM) -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${TAG} -t ${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) --build-arg VERSION=${VERSION} --push .
	@echo "-- Generate manifests for multiarch support --"
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:${TAG}" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${TAG}"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:${TAG}"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-${TAG}" --os=$(GOOS) --arch=$(GOARCH)
	docker manifest create --amend "${DOCKER_REGISTRY}/rabbitmq-client:latest" "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest"
	docker manifest annotate "${DOCKER_REGISTRY}/rabbitmq-client:latest"  "${DOCKER_REGISTRY}/rabbitmq-client:$(GOOS)-$(GOARCH)-latest" --os=$(GOOS) --arch=$(GOARCH)

manifests: ## Publish manifests for docker multiarch support
	@echo "-> $@"
	docker manifest push "${DOCKER_REGISTRY}/rabbitmq-client:${TAG}"
	docker manifest push "${DOCKER_REGISTRY}/rabbitmq-client:latest"



# -- build development binary --
.PHONY: build.darwin
build.darwin:
	@echo "-> $@"
	mkdir -p target/darwin
	go build -o target/darwin/mq_test -tags dev  it/main.go
	chmod 755 target/darwin/mq_test

.PHONY: test
test:
	@echo "-> $@"
	go test -tags=rabbitintegration  it/mq_test.gp -v -count=1;

.PHONY: fmt
fmt:
	@echo "-> $@"
	gofmt -s -l ./ | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint:
	@echo "-> $@"
	@go get -u golang.org/x/lint/golint
	@golint ./... | tee /dev/stderr
	@go get -u golang.org/x/tools/go/analysis/cmd/vet
	@go vet --all

.PHONY: dep
dep:
	@echo "-> $@"
	@go mod tidy
	@go get -u github.com/golang/dep/cmd/dep
	@dep init && dep ensure -vendor-only

.PHONY: imp
imp:
	@echo "-> $@"
	@goimport




