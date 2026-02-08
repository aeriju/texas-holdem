# Holdem Backend - AMD64 for GKE
# Build from project root: docker build -f backend/Dockerfile .
FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /app

COPY backend/go.mod ./backend/go.mod
RUN cd backend && go mod download 2>/dev/null || true

COPY backend/ ./backend/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./backend/cmd/api

# Runtime
FROM --platform=linux/amd64 alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /server .

EXPOSE 8080

ENTRYPOINT ["./server"]
