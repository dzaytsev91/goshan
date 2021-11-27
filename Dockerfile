# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app
COPY . ./
RUN go build main.go

CMD [ "/app/main" ]
