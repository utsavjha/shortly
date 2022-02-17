# syntax=docker/dockerfile:1
FROM golang:1.17-alpine


WORKDIR /app
COPY data* ./data
COPY db_clients* ./db_clients
COPY shortner_mod* ./shortner_mod

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /shortly

EXPOSE 8080

CMD [ "/shortly" ]