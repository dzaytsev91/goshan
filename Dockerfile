# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app
COPY go.mod ./
COPY main.go ./
COPY form.html ./
RUN go build main.go

CMD [ "/app/main" ]
