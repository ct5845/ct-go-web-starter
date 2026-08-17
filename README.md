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

4. Enable the repository git hooks (runs `gofmt` before each commit):
   ```bash
   git config core.hooksPath .githooks
   ```

5. Run the web development server:
   ```bash
   make web
   ```

The application will be available at `http://localhost:8080` (or the port set in `PORT`).

> `make web` detects the OS automatically — it uses `.air.windows.toml` on Windows and `.air.linux.toml` on Linux/macOS.

## Development

### Project Structure

```
├── cmd/
│   ├── web/           # Main entrypoint (starts the server, wires up routes)
│   └── copyassets/    # Build tool: copies static assets and JS deps to tmp/
├── src/
│   ├── features/      # Features with HTTP surface (routes + handlers)
│   │   ├── home/      # Home page feature
│   │   │   ├── home.go    # Handler, routes, and page assembly
│   │   │   └── home.html  # Feature template
│   │   ├── showcase/  # One page per component, for building them in isolation
│   │   └── nav/       # Shared navigation (bottom tabs, sidebar, more sheet)
│   ├── components/    # UI building blocks with no HTTP surface
│   │   ├── component/ # Component engine (New, Render, WithAlpine, WithIIFE)
│   │   ├── page/      # Full page shell template
│   │   ├── layoutswitch/ # Sidebar on desktop, bottom tabs on mobile
│   │   ├── layoutfull/   # Full-bleed layout with no chrome
│   │   ├── sidebar/      # Desktop navigation sidebar
│   │   ├── bottomtabs/   # Mobile bottom tab bar
│   │   ├── bottomsheet/  # Modal bottom sheet dialog
│   │   ├── pagedlist/    # Paginated list with htmx-driven search
│   │   └── icon/         # Icon font subsetting and inline SVGs
│   ├── infrastructure/ # Platform and runtime concerns
│   │   ├── config/    # Configuration and logging
│   │   ├── compression/ # HTTP response compression
│   │   └── fileserver/ # Static file serving with caching
│   └── static/        # Static assets (favicon, images, etc.)
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
3. Expose `RegisterRoutes(mux *http.ServeMux)` from the feature, and call it from `routes()` in `cmd/web/main.go`
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

### Showcase

`/showcase` renders every component on its own page, isolated from any feature. Use it to build and design a component before wiring it into a real page — each page names the file it comes from, and the slug matches the directory under `src/components/` or the stylesheet under `src/static/styles/`.

To add a page, append an entry to `demos` in [src/features/showcase/showcase.go](src/features/showcase/showcase.go) and write its render function. The index and the per-page navigation are both generated from that slice, so there is nothing else to update.

### Styling

TailwindCSS classes are available throughout the application. Modify `src/static/styles/styles.css` to add custom styles.

Tailwind scans `.html`, `.js`, and `.go` files, so class names written in Go — as the showcase's colour swatches are — are picked up too.

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
