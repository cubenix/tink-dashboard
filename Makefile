GO_FLAGS   ?=
VERSION    := v0.0.1
IMG_NAME   := quickdevnotes/tinker
IMG    	   := ${IMG_NAME}:${VERSION}

default: help

build: ## Build the CLI
	@go build ${GO_FLAGS} -o tinker main.go

img: ## Build the Docker Image
	@docker build --rm -t ${IMG} .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
