FROM alpine:3.11

ENTRYPOINT ["server"]
EXPOSE 7676

COPY src/app/templates /src/app/templates
COPY src/app/public /src/app/public

COPY server /bin/
