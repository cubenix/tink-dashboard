# -----------------------------------------------------------------------------
# The base image for building the tinker binary
FROM golang:1.17-alpine3.14 AS build

WORKDIR /tinker
COPY go.mod go.sum main.go Makefile ./
COPY internal internal
COPY cmd cmd
RUN apk --no-cache add make git gcc libc-dev curl && make build

# -----------------------------------------------------------------------------
# Build the final Docker image

FROM alpine:3.14.3

COPY --from=build /tinker/tinker /bin/tinker
ENTRYPOINT [ "/bin/tinker" ]
