FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o /urls-server

# TODO: investigate if scratch can be better (linking issues because of sqlite?)
FROM debian:trixie-slim as production
ENV GIN_MODE=release
WORKDIR /
COPY --from=builder /urls-server /urls-server
COPY ./migrations /migrations
VOLUME [ "/db" ]
EXPOSE 8080
ENTRYPOINT [ "/urls-server" ]
CMD [ "--migrate", "--migrations-folder=/migrations", "--sqlite-db-path=/db/urls.db" ]