
# Image URL to use all building/pushing image targets

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run go test again code.
	go test ./... -coverprofile cover.out

.PHONY: report
report: fmt vet ## Run go test again code.
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html

##@ Build

.PHONY: build
build: fmt vet ## Build binary.
	go build -o ltsctl main.go

.PHONY: run-help
run-help: fmt vet ## Display help how to run a ltsctl from your host.
	go run ./main.go 

.PHONY: run-wms-agent
run-wms-agent: fmt vet ## run a wms-agent using ltsctl from your host.
	go run ./main.go run wms-agent 

.PHONY: run-uploader
run-uploader: fmt vet ## run a uploader using ltsctl from your host.
	go run ./main.go run uploader

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: test ## Build docker image 
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image 
	docker push ${IMG}

.PHONY: docker-login
docker-login: ## Login docker image 
	aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}

.PHONY: docker-build-push
docker-build-push: docker-build docker-login docker-push## Build, Login and Push docker image 
	@echo "docker build and push Done"

.PHONY: docker-build-kwos
docker-build-kwos: test ## Build docker image for kwos
	docker build -t ${IMG_KWOS} -f ${DOCKERFILE_KWOS} .

.PHONY: docker-push-kwos
docker-push-kwos: ## Push docker image for kwos
	docker push ${IMG_KWOS}

.PHONY: docker-build-push-kwos
docker-build-push-kwos: docker-build-kwos docker-login docker-push-kwos ## Build, Login and Push docker image for kwos
	@echo "docker build and push Done"

.PHONY: docker-build-debug
docker-build-debug: test ## Build docker image for debug
	docker build -t ${IMG_DEBUG} -f ${DOCKERFILE_DEBUG} .

.PHONY: docker-push-debug
docker-push-debug: ## Push docker image for debug
	docker push ${IMG_DEBUG}

.PHONY: docker-build-push-debug
docker-build-push-debug: docker-build-debug docker-login docker-push-debug ## Build, Login and Push docker image for debug
	@echo "docker build and push Done"

.PHONY: docker-build-push-all
docker-build-push-all: docker-build docker-build-kwos docker-login docker-push docker-push-kwos ## Build, Login and Push docker image for all
	@echo "docker build and push Done"
