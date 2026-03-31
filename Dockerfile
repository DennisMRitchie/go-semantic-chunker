FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/semantic-chunker .

# ── Final slim image ──────────────────────────────────────
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/semantic-chunker .

EXPOSE 8083
CMD ["./semantic-chunker"]
