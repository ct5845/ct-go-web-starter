# Stage 1: build the Go binary
FROM golang:1.26.4-alpine AS builder
ARG VERSION=dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X ct-go-web-starter/internal/service.Version=${VERSION}" \
    -o mcp ./cmd/mcp

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static-debian12
WORKDIR /app
# Structured JSON logs on stdout, which is what cluster log collectors expect.
ENV APP_ENV=prod
# Each binary is alone in its container, so they all scrape on one port.
ENV METRICS_PORT=9090
COPY --from=builder /app/mcp .
EXPOSE 8082
# Scrape port. Keep it out of the Service — /metrics exposes internals.
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/app/mcp"]
