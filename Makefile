PROJECT_NAME	:= $(shell basename "$(PWD)")
GO_ENV 		  	:= CGO_ENABLED=0 GOOS=linux
BASE_DIR 		:= /src/server 
TLS_DIR 		:= $(shell echo `pwd`/tls)
TAG    			:= $(shell git log -1 --pretty=%H)
IMG    			:= ${PROJECT_NAME}:${TAG}
REDIS_DIR		:= $(shell echo `pwd`/redis-data)

export GO111MODULE=on
include .env

default: build

build:
	$(GO_ENV) go build -o server .$(BASE_DIR)
	docker build -t ${IMG} .
	docker tag ${IMG} ${PROJECT_NAME}:latest
	@echo "Build completed successfully"

certs:
	mkdir -p ${TLS_DIR}/certs
	docker build -t certs-generator ./tls
	docker run -t -v $(TLS_DIR)/certs:/certs certs-generator

redis:
	mkdir -p $(REDIS_DIR)
	docker run -d \
	-v $(REDIS_DIR):/data \
	-p 6379:6379 \
	redis:alpine redis-server --appendonly yes

run:
	docker run -d \
	--network ${TINKERBELL_NETWORK} \
	-e ALLOW_INSECURE=${ALLOW_INSECURE} \
	-e TINKERBELL_GRPC_AUTHORITY=${TINKERBELL_HOST}:42113 \
	-e TINKERBELL_CERT_URL=http://${TINKERBELL_HOST}:42114/cert \
	-e REDIS_ADDRESS=${REDIS_ADDRESS} \
	-v $(TLS_DIR)/certs:/tls \
	-p 7676:7676 \
	${PROJECT_NAME}:latest

	@if [ ${ALLOW_INSECURE} == true ] ; \
	then \
	echo "server listening at http://localhost:7676" ; \
	else \
	echo "server listening at https://localhost:7676" ; \
	fi ;

clean:
	rm -rf ${TLS_DIR}/certs
	rm server
