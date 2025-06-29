# Set the default task
default:
  just -l

# Allow APP_NAME to be overridden in the environment
APP_NAME := "./bootstrap"

# Build for linux/amd64 with CGO disabled
build:
  GOOS=linux GOARCH=amd64 go build -o {{APP_NAME}} ./cmd/main.go

# Alias `just dev` to `just run`
alias dev := run

# Run in dev mode using air
run:
  DEV=true air

wasm:
  GOOS=js GOARCH=wasm go build -o ./cmd/wasm/main.wasm ./cmd/wasm

# Static analysis
lint:
  staticcheck ./...

# Run tests with coverage
test:
  go test ./...

# Vet the code
vet:
  go vet ./...

# Format the code
fmt:
  go fmt ./...

# Clean up the binary
clean:
  rm {{APP_NAME}}
