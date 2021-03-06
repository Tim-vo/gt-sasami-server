# EXECUTABLE := gt-sasami-server.exe
# GITVERSION := $(shell git describe --dirty --always --tags --long)
# GOPATH ?= ${HOME}/go
# PACKAGENAME := $(shell go list -m -f '{{.Path}}')
# TOOLS := ${GOPATH}/bin/mockery \
# 	${GOPATH}/bin/swag
# SWAGGERSOURCE = $(wildcard gt-sasami-server/*.go) \
# 	$(wildcard gt-sasami-server/rpc/*.go)

# .PHONY: default
# default: ${EXECUTABLE}

# tools: ${TOOLS}

# ${GOPATH}/bin/mockery:
# 	go install github.com/vektra/mockery/cmd/mockery@latest

# ${GOPATH}/bin/swag:
# 	go install github.com/swaggo/swag/cmd/swag@latest

# .PHONY: swagger
# swagger: tools ${SWAGGERSOURCE}
# 	swag init --dir . --generalInfo gt-sasami-server/swagger.go --exclude embed --output embed/public_html/apidocs
# 	rm embed/public_html/apidocs/docs.go

# embed/public_html/apidocs/swagger.json: tools ${SWAGGERSOURCE}
# 	swag init --dir . --generalInfo gt-sasami-server/swagger.go --exclude embed --output embed/public_html/apidocs
# 	rm embed/public_html/apidocs/docs.go

# .PHONY: mocks
# mocks: tools
# 	mockery -dir ./gt-sasami-server -name GTStore

# .PHONY: ${EXECUTABLE}
# ${EXECUTABLE}: tools embed/public_html/apidocs/swagger.json
# 	# Compiling...
# 	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

# .PHONY: test
# test: tools mocks
# 	go test -cover ./...

# .PHONY: deps
# deps:
# 	# Fetching dependencies...
# 	go get -d -v # Adding -u here will break CI

# .PHONY: lint
# lint:
# 	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v --timeout 5m

# .PHONY: hadolint
# hadolint:
# 	docker run -it --rm -v ${PWD}/Dockerfile:/Dockerfile hadolint/hadolint:latest hadolint --ignore DL3018 Dockerfile

# .PHONY: relocate
# relocate:
# 	@test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
# 	$(eval ESCAPED_PACKAGENAME := $(shell echo "${PACKAGENAME}" | sed -e 's/[\/&]/\\&/g'))
# 	$(eval ESCAPED_TARGET := $(shell echo "${TARGET}" | sed -e 's/[\/&]/\\&/g'))
# 	# Renaming package ${PACKAGENAME} to ${TARGET}
# 	@grep -rlI '${PACKAGENAME}' * | xargs -i@ sed -i 's/${ESCAPED_PACKAGENAME}/${ESCAPED_TARGET}/g' @
# 	# Complete... 
# 	# NOTE: This does not update the git config nor will it update any imports of the root directory of this project.
GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/gt-sasami-server

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web