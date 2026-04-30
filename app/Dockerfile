FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o notiflex-api .

FROM scratch
COPY --from=builder /app/notiflex-api /notiflex-api
EXPOSE 8080
ENTRYPOINT ["/notiflex-api"]
