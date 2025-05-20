ARG PORT

FROM docker.io/library/golang:1.24.2-alpine3.20 AS build-stage

WORKDIR /app

COPY . .

RUN cd ./app_receiver && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o receiver

FROM docker.io/library/alpine:3.20

WORKDIR /app

COPY --from=build-stage /app/app_receiver/receiver receiver

EXPOSE $PORT

ENTRYPOINT ["/app/receiver"]
