# syntax=docker/dockerfile:1
FROM golang:1.17-alpine as build


WORKDIR /app
COPY data* ./data
COPY db_clients* ./db_clients
COPY shortner_mod* ./shortner_mod
COPY workers* ./workers

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /shortly

## Deploy Stage
FROM alpine:latest as runner

WORKDIR /

COPY --from=build /shortly /shortly

EXPOSE 8080

ENTRYPOINT [ "/shortly" ]