# Stage 1: build CSS and copy frontend assets
FROM node:24.16-alpine AS assets
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
COPY src/static ./src/static
RUN npx tailwindcss -i ./src/static/styles/styles.css -o ./tmp/static/style.css --minify
RUN cp node_modules/alpinejs/dist/cdn.min.js ./tmp/static/alpine.min.js
RUN cp node_modules/htmx.org/dist/htmx.min.js ./tmp/static/htmx.min.js
COPY src/static ./tmp/static

# Stage 2: build the Go binary
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o web ./cmd/web

# Stage 3: minimal runtime image
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/web .
COPY --from=assets /app/tmp/static ./tmp/static
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/web"]
