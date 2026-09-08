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
- **Three Binaries, One Service** - Web, JSON API, and MCP server over the same `internal/service`
- **MCP Server** - Streamable HTTP, added to LLM clients as a remote connector by URL

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
│   ├── web/           # Web app entrypoint (htmx + Alpine, port 8080)
│   ├── api/           # JSON API entrypoint (port 8081)
│   └── mcp/           # MCP server entrypoint (port 8082)
├── tools/
│   └── copyassets.mjs # Build step: static assets and npm libs into tmp/
├── internal/
│   ├── service/       # Domain types and operations — imports nothing else here
│   ├── api/           # JSON handlers over the service
│   ├── mcpserver/     # MCP tool definitions over the service
│   ├── web/           # Everything cmd/web needs
│   │   ├── features/  # Features with HTTP surface (routes + handlers)
│   │   │   ├── home/         # Home page feature
│   │   │   ├── showcase/     # One page per component, built in isolation
│   │   │   ├── mcpconnector/ # Connector page: endpoint URL and client config
│   │   │   ├── notfound/     # Catch-all 404 page
│   │   │   └── nav/          # Shared navigation (tabs, sidebar, more sheet)
│   │   ├── components/# UI building blocks with no HTTP surface
│   │   │   ├── component/    # Component engine (New, Render, WithAlpine)
│   │   │   ├── demo/         # Showcase page type, declared by each component
│   │   │   ├── page/         # Full page shell template
│   │   │   ├── layoutswitch/ # Sidebar on desktop, bottom tabs on mobile
│   │   │   └── icon/         # Icon font subsetting and inline SVGs
│   │   └── static/    # Static assets (favicon, images, stylesheets)
│   └── infrastructure/# Platform and runtime concerns
│       ├── config/      # Configuration and logging
│       ├── health/      # /healthz, mounted by every HTTP binary
│       ├── httpserver/  # Listen and graceful shutdown
│       ├── compression/ # HTTP response compression
│       └── fileserver/  # Static file serving with caching
├── build/             # Production binary output (not in git)
├── tmp/               # Dev build output (not in git)
├── .air.windows.toml  # Live reload config (Windows)
├── .air.linux.toml    # Live reload config (Linux/macOS/devcontainer)
└── package.json       # Frontend dependencies
```

### Available Commands

- `make web` - Start web development server with live reload
- `make up` / `make down` - Run all three services together in Docker
- `make build` - Build CSS, copy assets, and compile everything under `cmd/`
- `make test` - Run `go vet` and the test suite
- `make docker` - Build all three production Docker images

### Running Everything

`make web` is the daily loop — native, live-reloading, sub-second rebuilds. The
web app makes no outbound calls, so nothing else needs to be running for it.

To point an LLM client at the MCP server while developing, run it alongside in a
second terminal; `MCP_PUBLIC_URL` already defaults to where it lands, so the
config on `/connector` is pasteable as-is:

```bash
go run ./cmd/mcp
```

To see all three running together as they are built and shipped:

```bash
make up     # builds the images and starts the stack
make down
```

That is a check that everything works together, not a development loop — anything
containerised costs an image rebuild per change. Keep editing under `make web`.

### The Three Binaries

All three bind `internal/service` directly — the Go interface is the contract, so
the web app never calls its own API over HTTP.

| Binary | Default port | Surface |
| --- | --- | --- |
| `cmd/web` | 8080 | Server-rendered pages, htmx + Alpine |
| `cmd/api` | 8081 | JSON, e.g. `GET /v1/greet?name=Ada` |
| `cmd/mcp` | 8082 | MCP over streamable HTTP at `/mcp` |

Each listens on `PORT`; the defaults differ so all three run side by side
locally. Each answers two probes, both carrying the build version stamped in via
ldflags:

- `GET /healthz` — liveness. Checks no dependencies on purpose: a failing
  database should take pods out of service, not restart every one of them.
- `GET /ready` — readiness. Returns 503 the moment shutdown begins, so traffic
  stops being routed before the listener closes.

Neither appears in the request log, so probes don't drown it — but both are
counted in metrics, because a probe starting to fail is exactly what you want to
see.

Each binary also serves Prometheus metrics on a **separate port** (`METRICS_PORT`,
default 9090): request counts and latency histograms labelled by route, plus Go
runtime and process metrics. The separate port is deliberate — `/metrics`
exposes internals and must never ride on the port that serves users.

Metrics are labelled with the matched **route pattern**, never the request path.
`GET /showcase/{slug}` stays one series no matter how many slugs are requested;
labelling by path would create a series per URL and eventually take the scraper
down.

Domain errors are typed values on the service and each transport translates
them — `service.ErrEmptyName` becomes a 400 in the API and a tool-level error in
MCP, never a 500.

### MCP Connector

`cmd/mcp` serves the tools over streamable HTTP at `/mcp`. Clients add it as a
remote connector by URL — there is nothing to download or install, and the
credentials the tools need stay on the server.

The web app's `/connector` page shows that URL and a ready-made client config.
Point it at the right host by setting `MCP_PUBLIC_URL`.

### Adding New Features

1. Create a new feature directory in `internal/web/features/`
2. Add a `.go` file with routes, handler, and page assembly
3. Expose `RegisterRoutes(mux *http.ServeMux)` from the feature, and call it from `routes()` in `cmd/web/main.go`
4. Use components from `internal/web/components/` or create feature-internal ones in the feature directory

**Example: Adding a "blog" feature**
```
internal/web/features/blog/
├── blog.go        # Routes, handler, and page assembly
├── list.html      # Blog listing template
├── postcard.go    # Feature-internal component (unexported)
└── postcard.html
```

Split into `handler.go` + `page.go` only if page assembly grows complex enough to warrant it.

### Showcase

`/showcase` renders every component on its own page, isolated from any feature. Use it to build and design a component before wiring it into a real page — each page names the file it comes from, and the slug matches the directory under `internal/web/components/` or the stylesheet under `internal/web/static/styles/`.

Each component's demo lives in the component's own directory as `showcase.go`, exporting `var Showcase = demo.Page{...}` — so the component and its demo show up in the same diff. The feature holds only the routes, the index, and the page chrome, plus the demos for stylesheets, which have no component directory to live in.

To add a page: write `showcase.go` next to the component, then add its `Showcase` value to `demos` in [internal/web/features/showcase/showcase.go](internal/web/features/showcase/showcase.go). The index and the per-page navigation are both generated from that slice.

### Styling

TailwindCSS classes are available throughout the application. Modify `internal/web/static/styles/styles.css` to add custom styles.

The frontend build is a single npm script, `build-frontend`, used everywhere —
by `make web` (via air), by `make build`, and by the web Dockerfile. It compiles
the stylesheet and then copies the static tree and the npm-installed libraries
into `tmp/static`. Adding a frontend dependency means adding one line to
`tools/copyassets.mjs`; there is no second place to update.

Tailwind scans `.html`, `.js`, and `.go` files, so class names written in Go — as the showcase's colour swatches are — are picked up too.

## Dev Container

A [Dev Container](https://containers.dev/) configuration is included, giving
every developer the same Linux toolchain regardless of host OS. It needs Docker
Desktop and the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).

`F1` → **"Dev Containers: Reopen in Container"**. The first build takes a few
minutes; after that it opens in seconds. `F1` → "Reopen Folder Locally" gets you
back out.

You get Go 1.26, Node 24, Air, `make`, and a Docker CLI wired to the host daemon,
so every `make` target works inside — including `make up` and `make docker`. All
six ports are forwarded; 8080 opens a browser and the rest are silent. AWS
credentials are passed through from host environment variables, not mounted.

### Why use it

The native toolchain is faster, so if you already have Go and Node set up
locally, keep using them for day-to-day work. The dev container earns its place
in three situations:

- **Onboarding.** Clone, open, work. No one spends a morning matching Go and
  Node versions.
- **Reproducing Linux-only behaviour.** Production is Linux containers. Some
  things simply do not reproduce on Windows or macOS — signal handling during
  graceful shutdown, file ownership and permissions, case-sensitive paths.
- **Pinning versions across a team**, so "works on my machine" stops being a
  category of bug.

### node_modules is not shared with the host

It lives in a named Docker volume rather than the bind-mounted workspace.
Tailwind and Parcel's file watcher ship platform-native binaries, so a shared
`node_modules` means whichever platform ran `npm install` last breaks the other
— with errors like `No prebuild or local build of @parcel/watcher found`.

The volume keeps a Linux install inside and the host's install outside, and
neither disturbs the other. If you ever do see that error on the host, run
`npm install` to repair it.

`tmp/` and `build/` are still shared, but hold platform-specific binaries
(`build/web` vs `build/web.exe`). They coexist harmlessly; if anything looks odd
after switching, delete both and rebuild. Nothing there is tracked.

## Production

Build the application for production:

```bash
make build
./build/web
```

`make build` produces `web`, `api` and `mcp` in `build/`. Each binds to
`0.0.0.0:<PORT>`, so they work in containers and behind reverse proxies. Set
`PORT` via environment variable — no `.env` file is required in production.

### Kubernetes

The images are built for a cluster: distroless with no shell or package manager,
`USER nonroot`, bound to all interfaces, configured entirely by environment, and
logging JSON to stdout (`APP_ENV=prod` is set in the image).

On `SIGTERM` each binary stops reporting ready, keeps serving for 5 seconds
while endpoints propagate, then drains in-flight requests — up to 10 more
seconds. **`terminationGracePeriodSeconds` must stay above 15**; the default of
30 is fine, but don't lower it.

```yaml
spec:
  terminationGracePeriodSeconds: 30
  containers:
    - name: web
      image: ct-go-web-starter-web:v0.1.0
      ports:
        - name: http
          containerPort: 8080
        # Scrape target. Deliberately not in the Service — /metrics exposes
        # internals and must not be reachable from where users are.
        - name: metrics
          containerPort: 9090
      env:
        - name: PORT
          value: "8080"
        # Go reads cgroup CPU limits by itself, but not memory limits. Without
        # this the GC can let the heap run past the limit and be OOMKilled.
        - name: GOMEMLIMIT
          valueFrom:
            resourceFieldRef:
              resource: limits.memory
      readinessProbe:
        httpGet: { path: /ready, port: 8080 }
        periodSeconds: 5
        failureThreshold: 2
      livenessProbe:
        httpGet: { path: /healthz, port: 8080 }
        periodSeconds: 15
        failureThreshold: 3
      securityContext:
        runAsNonRoot: true
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
        capabilities: { drop: ["ALL"] }
      resources:
        requests: { cpu: 100m, memory: 64Mi }
        limits: { cpu: "1", memory: 256Mi }
```

Two things worth knowing. Go 1.25+ derives `GOMAXPROCS` from the cgroup CPU
**limit** — not the request — so a pod with no limit set will assume every core
on the node. And `cmd/mcp` runs its MCP transport in stateless mode, so replicas
need no session affinity and scale horizontally behind a plain Service.

No startup probe is needed; these bind in milliseconds.

### Docker

Build the image:

```bash
make docker
```

This builds three images: `ct-go-web-starter-web`, `-api` and `-mcp`, each from
its own Dockerfile under `cmd/`.

Run one:

```bash
docker run -p 8080:8080 ct-go-web-starter-web
```

Then open `http://localhost:8080`. `Ctrl+C` to stop — the server drains in-flight requests before exiting.

To run on a different port:

```bash
docker run -p 9000:9000 -e PORT=9000 ct-go-web-starter-web
```

The web image uses a three-stage build (Node → Go → distroless); the API and MCP
images skip the Node stage. All three produce a minimal runtime image with no shell or package manager. Environment variables are the only supported configuration mechanism — no `.env` file is used at runtime.

## License

This project is open source and available under the [MIT License](LICENSE).
