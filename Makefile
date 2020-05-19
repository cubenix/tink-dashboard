PROJECT_NAME	:= $(shell basename "$(PWD)")
GO_ENV 		  	:= CGO_ENABLED=0 GOOS=linux
BASE_DIR 		:= /src/server 
TLS_DIR 		:= $(shell echo `pwd`/tls)
TAG    			:= $(shell git log -1 --pretty=%H)
IMG    			:= ${PROJECT_NAME}:${TAG}
REDIS_DIR		:= $(shell echo `pwd`/redis-data)

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

redis:
	mkdir -p $(REDIS_DIR)
	docker run -d \
		--name twiz-redis \
		-v $(REDIS_DIR):/data \
		-p 6379:6379 \
		redis:alpine redis-server --appendonly yes

run: build gen-certs
	docker run -d \
		--name ${PROJECT_NAME} \
		-e ALLOW_INSECURE=false \
		-e TINKERBELL_GRPC_AUTHORITY=${TINKERBELL_HOST_IP}:42113 \
      	-e TINKERBELL_CERT_URL=http://${TINKERBELL_HOST_IP}:42114/cert \
		-p 7676:7676 \
		${PROJECT_NAME}:latest
	echo "server listening at https://localhost:7676"

run-insecure: build 
	docker run -d \
		--name ${PROJECT_NAME} \
		--network ${TINKERBELL_NETWORK} \
		-e ALLOW_INSECURE=true \
		-e TINKERBELL_GRPC_AUTHORITY=${TINKERBELL_HOST}:42113 \
      	-e TINKERBELL_CERT_URL=http://${TINKERBELL_HOST}:42114/cert \
		-e REDIS_ADDRESS=192.168.1.4:6379 \
		-p 7676:7676 \
		${PROJECT_NAME}:latest
	echo "server listening at http://localhost:7676"
	docker logs -f tink-wizard

clean:
	rm -rf ${TLS_DIR}/certs
	rm server

all: build-server
