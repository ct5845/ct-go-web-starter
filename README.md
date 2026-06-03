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

- Go 1.24.2 or later
- Node.js (for TailwindCSS and frontend dependencies)

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

4. Run the development server:
   ```bash
   make dev
   ```

The application will be available at `http://localhost:8080` (or the port set in `PORT`).

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
├── .air.toml          # Live reload configuration
└── package.json       # Frontend dependencies
```

### Available Commands

- `make dev` - Start development server with live reload (runs air)
- `make build` - Build CSS, copy assets, and compile production binary

### Adding New Features

1. Create a new feature directory in `src/features/`
2. Add `handler.go` and feature-specific templates
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

TailwindCSS classes are available throughout the application. Modify `src/styles/styles.css` to add custom styles or extend the configuration in `src/styles/config.css`.

## Production

Build the application for production:

```bash
make build
./build/web
```

The server binds to `0.0.0.0:<PORT>` (default `8080`), so it works in containers and behind reverse proxies. Set `PORT` via environment variable — no `.env` file is required in production.

## License

This project is open source and available under the [MIT License](LICENSE).