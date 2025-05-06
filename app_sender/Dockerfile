FROM docker.io/library/golang:1.24.2-alpine3.20 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sender

FROM docker.io/library/alpine:3.20

WORKDIR /app

COPY --from=build-stage /app/sender sender

EXPOSE 54321

ENTRYPOINT ["/app/sender"]
