FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o notiflex .

FROM scratch
COPY --from=builder /app/notiflex /notiflex
EXPOSE 8080
ENTRYPOINT ["/notiflex"]
