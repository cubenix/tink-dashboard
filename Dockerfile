FROM alpine:3.11

ENTRYPOINT ["server"]
EXPOSE 7676

RUN apk add --update ca-certificates && \
    apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing cfssl

COPY server /bin/
COPY src/app/templates /src/app/templates
COPY src/app/public /src/app/public
