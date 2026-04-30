# Stage 1: Build the Go binary
FROM golang:1.25-alpine as builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /notiflex-api

# Stage 2: Create the final small image
FROM scratch

COPY --from=builder /notiflex-api /notiflex-api

EXPOSE 8080

CMD ["/notiflex-api"]
