# CT Go Web Starter

A modern Go web application starter template with HTMX, Alpine.js, and TailwindCSS.

## Features

- **Go** - Clean, fast backend with structured routing
- **HTMX** - Dynamic frontend interactions without JavaScript complexity
- **Alpine.js** - Lightweight JavaScript framework for reactivity
- **TailwindCSS** - Utility-first CSS framework
- **Live Reload** - Air integration for development hot reloading
- **Static Asset Caching** - Built-in ETag support for efficient caching
- **Gzip Compression** - Automatic response compression for supporting clients
- **Graceful Shutdown** - Drains in-flight requests on SIGINT/SIGTERM
- **Feature-Based Architecture** - Organized by features (vertical slices) for better maintainability

## Quick Start

### Prerequisites

- Go 1.26 or later
- Node.js 24 or later
- [Air](https://github.com/air-verse/air) for live reload (`go install github.com/air-verse/air@latest`)

### Installation

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd ct-go-web-starter
   ```

2. Install dependencies:
   ```bash
   npm install
   go mod tidy
   ```

3. Copy `.env.example` to `.env` and adjust as needed:
   ```bash
   cp .env.example .env
   ```

4. Run the web development server:
   ```bash
   make web
   ```

The application will be available at `http://localhost:8080` (or the port set in `PORT`).

> `make web` detects the OS automatically — it uses `.air.windows.toml` on Windows and `.air.linux.toml` on Linux/macOS.

## Development

### Project Structure

```
├── cmd/
│   ├── web/           # Main entrypoint (starts the server)
│   └── copyassets/    # Build tool: copies static assets and JS deps to tmp/
├── src/
│   ├── features/      # Features with HTTP surface (routes + handlers)
│   │   └── home/      # Home page feature
│   │       ├── home.go    # Handler, routes, and page assembly
│   │       └── home.html  # Feature template
│   ├── components/    # UI building blocks with no HTTP surface
│   │   ├── component/ # Component engine (New, Render, WithJS)
│   │   └── page/      # Full page shell template
│   ├── infrastructure/ # Platform and runtime concerns
│   │   ├── config/    # Configuration and logging
│   │   ├── compression/ # HTTP response compression
│   │   └── fileserver/ # Static file serving with caching
│   ├── static/        # Static assets (favicon, images, etc.)
│   └── app.go         # Application setup and routing
├── build/             # Production binary output (not in git)
├── tmp/               # Dev build output (not in git)
├── .air.toml          # Live reload config (Windows)
├── .air.linux.toml    # Live reload config (Linux/macOS/devcontainer)
└── package.json       # Frontend dependencies
```

### Available Commands

- `make web` - Start web development server with live reload
- `make build` - Build CSS, copy assets, and compile production binary
- `make docker` - Build the production Docker image

### Adding New Features

1. Create a new feature directory in `src/features/`
2. Add a `.go` file with routes, handler, and page assembly
3. Register routes in `src/app.go`
4. Use components from `src/components/` or create feature-internal ones in the feature directory

**Example: Adding a "blog" feature**
```
src/features/blog/
├── blog.go        # Routes, handler, and page assembly
├── list.html      # Blog listing template
├── postcard.go    # Feature-internal component (unexported)
└── postcard.html
```

Split into `handler.go` + `page.go` only if page assembly grows complex enough to warrant it.

### Styling

TailwindCSS classes are available throughout the application. Modify `src/static/styles/styles.css` to add custom styles.

## Dev Container

This project includes a [Dev Container](https://containers.dev/) configuration for a consistent, isolated development environment. It requires the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) and Docker Desktop.

To use it: open the repo in VS Code and click **Reopen in Container** when prompted, or run `F1` → "Dev Containers: Reopen in Container".

The container includes Go 1.26, Node 24, Air, and all project dependencies pre-installed. AWS credentials are passed in from host environment variables — no credential files are mounted.

## Production

Build the application for production:

```bash
make build
./build/web
```

The server binds to `0.0.0.0:<PORT>` (default `8080`), so it works in containers and behind reverse proxies. Set `PORT` via environment variable — no `.env` file is required in production.

### Docker

Build the image:

```bash
make docker
```

Run it:

```bash
docker run -p 8080:8080 ct-go-web-starter
```

Then open `http://localhost:8080`. `Ctrl+C` to stop — the server drains in-flight requests before exiting.

To run on a different port:

```bash
docker run -p 9000:9000 -e PORT=9000 ct-go-web-starter
```

The image uses a three-stage build (Node → Go → distroless) producing a minimal runtime image with no shell or package manager. Environment variables are the only supported configuration mechanism — no `.env` file is used at runtime.

## License

This project is open source and available under the [MIT License](LICENSE).
