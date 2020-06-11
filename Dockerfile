FROM golang:1.14-buster as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY src src
RUN CGO_ENABLED=0 go build -o tink-wizard -ldflags="-s" ./src/server

FROM scratch
WORKDIR /app
COPY --from=builder /app/tink-wizard .
COPY --from=builder /app/src/app/templates src/app/templates
COPY --from=builder /app/src/app/public src/app/public
EXPOSE 7676
ENTRYPOINT ["/app/tink-wizard"]
