# Stage 1: build CSS and copy frontend assets
FROM node:24.16-alpine AS assets
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
# Tailwind scans .go, .html and .js for class names, so the whole of internal/
# must be present — not just the static directory.
COPY internal ./internal
COPY tools/copyassets.mjs ./tools/
RUN npm run build-frontend

# Stage 2: build the Go binary
FROM golang:1.26.4-alpine AS builder
ARG VERSION=dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X ct-go-web-starter/internal/service.Version=${VERSION}" \
    -o web ./cmd/web

# Stage 3: minimal runtime image
FROM gcr.io/distroless/static-debian12
WORKDIR /app
# Structured JSON logs on stdout, which is what cluster log collectors expect.
ENV APP_ENV=prod
# Each binary is alone in its container, so they all scrape on one port.
ENV METRICS_PORT=9090
COPY --from=builder /app/web .
COPY --from=assets /app/tmp/static ./tmp/static
EXPOSE 8080
# Scrape port. Keep it out of the Service — /metrics exposes internals.
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/app/web"]
