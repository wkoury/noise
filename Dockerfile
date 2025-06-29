# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the WASM binary
RUN GOOS=js GOARCH=wasm go build -o cmd/wasm/main.wasm ./cmd/wasm

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM scratch

WORKDIR /root/

# Copy the binary and static files
COPY --from=builder /app/server .
COPY --from=builder /app/cmd/wasm/main.wasm ./cmd/wasm/
COPY --from=builder /app/internal/server/views ./internal/server/views/
COPY --from=builder /app/internal/server/static ./internal/server/static/

EXPOSE 8080

CMD ["./server"]