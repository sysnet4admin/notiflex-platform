FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/notiflex-api .

FROM scratch
COPY --from=builder /out/notiflex-api /notiflex-api

EXPOSE 8080
ENTRYPOINT ["/notiflex-api"]
