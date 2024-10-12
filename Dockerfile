FROM golang:1.23 AS build-stage

WORKDIR /src/app

COPY go.mod go.sum ./
COPY main.go main.go
COPY internal ./internal
RUN go build -o server main.go

FROM debian:bookworm

WORKDIR /root/
COPY --from=build-stage /src/app ./app
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

CMD ["./app/server"]