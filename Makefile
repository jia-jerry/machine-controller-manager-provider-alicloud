# Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROVIDER_NAME       := SampleProvider
PROJECT_NAME        := gardener
BINARY_PATH         := bin/
IMAGE_REPOSITORY    := docker-repository-link-goes-here
IMAGE_TAG           := $(shell cat VERSION)

#########################################
# Rules for running helper scripts
#########################################

.PHONY: rename-provider
rename-provider:
	@./hack/rename-provider ${PROVIDER_NAME}

.PHONY: rename-project
rename-project:
	@./hack/rename-project ${PROJECT_NAME}

#########################################
# Rules for re-vendoring
#########################################

.PHONY: revendor
revendor:
	@env GO111MODULE=on go mod vendor -v
	@env GO111MODULE=on go mod tidy -v

.PHONY: update-dependencies
update-dependencies:
	@env GO111MODULE=on go get -u

#########################################
# Rules for build/release
#########################################

.PHONY: release
release: build-local build docker-image docker-push rename-binaries

.PHONY: build-local
build-local:
	go build \
	-v \
	-o ${BINARY_PATH}/cmi-plugin \
	app/controller/cmi-plugin.go

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-a \
	-v \
	-o ${BINARY_PATH}/rel/cmi-plugin \
	app/controller/cmi-plugin.go

.PHONY: docker-image
docker-image:
	@if [[ ! -f ${BINARY_PATH}/rel/cmi-plugin ]]; then echo "No binary found. Please run 'make build'"; false; fi
	@docker build -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) .

.PHONY: docker-push
docker-push:
	@if ! docker images $(IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(IMAGE_TAG); then echo "$(IMAGE_REPOSITORY) version $(IMAGE_TAG) is not yet built. Please run 'make docker-images'"; false; fi
	@gcloud docker -- push $(IMAGE_REPOSITORY):$(IMAGE_TAG)

.PHONY: rename-binaries
rename-binaries:
	@if [[ -f bin/cmi-plugin ]]; then cp bin/cmi-plugin cmi-plugin-darwin-amd64; fi
	@if [[ -f bin/rel/cmi-plugin ]]; then cp bin/rel/cmi-plugin cmi-plugin-linux-amd64; fi

.PHONY: clean
clean:
	@rm -rf bin/
	@rm -f *linux-amd64
	@rm -f *darwin-amd64
