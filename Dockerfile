FROM golang:1.15-alpine AS builder

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/go-redis .

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

# Enable release mode for Gin framework
ENV GIN_MODE=release

COPY --from=builder /app/out/go-redis /app/go-redis
COPY templates templates

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/go-redis"]