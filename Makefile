# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
PROJECT_NAME=tap-nginx-dyn-stream

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

build: verify_gopath
	go fmt $(APP_DIR_LIST)
	CGO_ENABLED=0 go install -tags netgo .
	mkdir -p application && cp -f $(GOBIN)/$(APP_NAME) ./application/$(APP_NAME)

docker_build: build_anywhere
	docker build -t tap-api-service .

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/$(PROJECT_NAME)
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/$(PROJECT_NAME)

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/$(PROJECT_NAME)/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp ./$(PROJECT_NAME) ./application/$(PROJECT_NAME)
	rm -Rf ./temp

test: verify_gopath
	CGO_ENABLED=0 go test -tags netgo --cover $(APP_DIR_LIST)

mock_update:
	$(GOBIN)/mockgen -source=utils.go -package=main -destination=utils_mock_test.go
	$(GOBIN)/mockgen -source=systemUtils.go -package=main -destination=systemUtils_mock_test.go

