FROM golang:1.17-alpine

WORKDIR /app
COPY . ./
RUN go build main.go

CMD [ "/app/main" ]
