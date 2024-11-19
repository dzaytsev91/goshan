FROM golang:1.16-buster as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-api main.go
#######################################
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-api .
COPY --from=builder /app/assets/* /app/assets/
COPY --from=builder /app/form.html .
COPY --from=builder /app/image.html .
EXPOSE 8080
CMD ["./go-api"]
