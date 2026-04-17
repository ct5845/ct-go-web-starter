# CT Go Web Starter

A modern Go web application starter template with HTMX, Alpine.js, and TailwindCSS.

## Features

- **Go** - Clean, fast backend with structured routing
- **HTMX** - Dynamic frontend interactions without JavaScript complexity
- **Alpine.js** - Lightweight JavaScript framework for reactivity
- **TailwindCSS** - Utility-first CSS framework
- **Live Reload** - Air integration for development hot reloading
- **Static Asset Caching** - Built-in ETag support for efficient caching
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

3. Build CSS and assets:
   ```bash
   go generate
   ```

4. Run the development server:
   ```bash
   air
   ```

The application will be available at `http://localhost:8080`

## Development

### Project Structure

```
├── src/
│   ├── features/           # Features with HTTP surface (routes + handlers)
│   │   └── home/          # Home page feature
│   │       ├── handler.go # HTTP handler
│   │       ├── page.go    # Page assembly (subcomponents, data prep)
│   │       └── home.html  # Feature template
│   ├── components/        # UI building blocks with no HTTP surface
│   │   ├── component/     # Component engine (New, Render, WithJS)
│   │   ├── header/        # Site header
│   │   ├── footer/        # Site footer
│   │   └── page/          # Full page shell template
│   ├── infrastructure/    # Platform and runtime concerns
│   │   ├── config/        # Configuration and logging
│   │   └── fileserver/    # Static file serving with caching
│   ├── static/           # Static assets (favicon, images, etc.)
│   ├── styles/           # TailwindCSS source files
│   └── app.go           # Application setup and routing
├── scripts/             # Build scripts
├── build/              # Generated assets (not in git)
├── .air.toml          # Live reload configuration
└── package.json       # Frontend dependencies
```

### Available Commands

- `air` - Start development server with live reload
- `go run main.go` - Run without live reload
- `go generate` - Build CSS and copy assets
- `npm run build-css` - Build TailwindCSS only

### Adding New Features

1. Create a new feature directory in `src/features/`
2. Add `handler.go` and feature-specific templates
3. Register routes in `src/app.go`
4. Use components from `src/components/` or create feature-internal ones in the feature directory

**Example: Adding a "blog" feature**
```
src/features/blog/
├── handler.go     # HTTP handler
├── page.go        # Page assembly
├── list.html      # Blog listing template
└── postcard.go    # Feature-internal component (unexported)
└── postcard.html
```

### Styling

TailwindCSS classes are available throughout the application. Modify `src/styles/styles.css` to add custom styles or extend the configuration in `src/styles/config.css`.

## Production

Build the application for production:

```bash
go generate
go build -o app main.go
./app
```

## License

This project is open source and available under the [MIT License](LICENSE).