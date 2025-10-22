# Commands for building application

ROOT_DIR=$(PWD)

build: # Build everything
	@make gen
	@make goimports
	@make build-containers IMAGE_TAG=$(IMAGE_TAG)

build-containers:
	@make build-frontend
	@make build-backend

build-frontend:
	@make web-docker

build-backend:
	@make app-docker optionalReproFlag=$(optionalReproFlag)
	@make ffmpeg-docker optionalReproFlag=$(optionalReproFlag)

# Base makefile target for building a binary
GOOS_OVERRIDE ?= GOOS=linux
build-binary:
	@echo "Building $(executablePath) with tag: $(IMAGE_TAG)"
	@cd $(directory) && \
	CGO_ENABLED=0 $(GOOS_OVERRIDE) $(extraArgs) go build -buildvcs=false -asmflags= -trimpath -ldflags "-buildid= -s -extldflags "-static"" && \
	cd $(ROOT_DIR) && \
	sha256sum $(executablePath)

# Base makefile target for building a docker image
build-docker:
	@DOCKER_BUILDKIT=1 docker build $(buildPlatform) \
	--target $(target) \
	--quiet \
	. \
	-t $(target):latest \
	--build-arg ALPINE_CONTAINER_IMAGE=$(ALPINE_CONTAINER_IMAGE) \
	--build-arg GO_CONTAINER_IMAGE=$(GO_CONTAINER_IMAGE) \
	$(additionalDockerArgs)

build-docker-ffmpeg:
	@DOCKER_BUILDKIT=1 docker build $(buildPlatform) \
	--file Dockerfile.ffmpeg \
	--target srl-ffmpeg \
	--quiet \
	. \
	-t srl-ffmpeg:latest \
	--build-arg ALPINE_CONTAINER_IMAGE=$(ALPINE_CONTAINER_IMAGE) \
	$(additionalDockerArgs)

build-docker-web:
	@DOCKER_BUILDKIT=1 docker build $(buildPlatform) \
	--file web/Dockerfile \
	--quiet \
	web \
	-t srl-web:latest \
	$(additionalDockerArgs)

app-build:
	@make build-binary extraArgs=$(extraArgs) directory=cmd/app executablePath=cmd/app/app

app-docker:
	@make app-build $(optionalReproFlag) extraArgs=$(extraArgs)
	@make build-docker buildPlatform=$(buildPlatorm) target=srl

ffmpeg-docker:
	@make app-build $(optionalReproFlag) extraArgs=$(extraArgs)
	@make build-docker-ffmpeg buildPlatform=$(buildPlatorm) target=srl-ffmpeg

web-docker:
	@make build-docker-web buildPlatform=$(buildPlatorm) target=srl-web
