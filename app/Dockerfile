# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /notiflex-api .

# Stage 2: Create a minimal image
FROM scratch
WORKDIR /
COPY --from=builder /notiflex-api /notiflex-api
USER 1000
ENTRYPOINT ["/notiflex-api"]
