PROJECT_NAME	:= $(shell basename "$(PWD)")
GO_ENV 		  	:= CGO_ENABLED=0 GOOS=linux
BASE_DIR 		:= /src/server 
TLS_DIR 		:= $(shell echo `pwd`/tls)
TAG    			:= $(shell git log -1 --pretty=%H)
IMG    			:= ${PROJECT_NAME}:${TAG}

export GO111MODULE=on
MAKEFLAGS += --silent

default: all

vet:
	$(GO_ENV) go vet `go list -f '{{.Dir}}' .$(BASE_DIR) `	

fmt: vet
	$(GO_ENV) gofmt -w -s `go list -f '{{.Dir}}' .$(BASE_DIR) `

install:
	$(GO_ENV) go install `go list -f '{{.Dir}}' .$(BASE_DIR) `

build-server: fmt install
	$(GO_ENV) go build -o server .$(BASE_DIR)

build: build-server
	docker build -t ${IMG} .
	docker tag ${IMG} ${PROJECT_NAME}:latest

gen-certs:
	mkdir -p ${TLS_DIR}/certs
	docker build -t certs-generator ./tls
	docker run -t -v $(TLS_DIR)/certs:/certs certs-generator

run: build gen-certs
	docker run -d -v ${TLS_DIR}/certs:/tls -p 7676:7676 ${PROJECT_NAME}:latest sh
	echo "server listening at https://localhost:7676"

run-insecure: build
	docker run -d -e ALLOW_INSECURE=true -p 7676:7676 ${PROJECT_NAME}:latest sh
	echo "server listening at http://localhost:7676"

clean:
	rm -rf ${TLS_DIR}/certs
	rm server

all: build-server
