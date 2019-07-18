# Copyright 2016 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DOCKER_ARCHS      := amd64
DOCKER_IMAGE_NAME := prometheus-mailgun-exporter
DOCKER_IMAGE_TAG  := $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKER_REPO       := missionlane

BUILD_DATE   := $(shell date +%Y%m%d-%H:%M:%S)
BUILD_USER   := $(shell whoami)
GIT_BRANCH   := $(shell git rev-parse --abbrev-ref HEAD)
GIT_REVISION := $(shell git rev-parse HEAD)
GO111MODULE  := on
VERSION      := $(shell cat VERSION)

.PHONY: all
all: prometheus-mailgun-exporter docker

.PHONY: fmt
fmt:
	@GO_FMT_RESULT=$$(gofmt -d .); \
	if [ -n "$${GO_FMT_RESULT}" ]; then \
		echo "gofmt checking failed!"; \
		echo "$${GO_FMT_RESULT}"; \
		exit 1; \
	fi

prometheus-mailgun-exporter: fmt
	@go build \
		-ldflags "-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE) \
			-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) \
			-X github.com/prometheus/common/version.Branch=$(GIT_BRANCH) \
			-X github.com/prometheus/common/version.Revision=$(GIT_REVISION) \
			-X github.com/prometheus/common/version.Version=$(VERSION)"

.PHONY: docker
docker: fmt
	@docker build \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GIT_REVISION=$(GIT_REVISION) \
		--build-arg GO111MODULE=$(GO111MODULE) \
		--build-arg VERSION=$(VERSION) \
		-t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" \
		.
