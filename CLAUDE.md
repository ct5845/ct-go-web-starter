# Coding Guidelines

These guidelines apply to all code in this project. They are intended to steer both human and AI-assisted development toward a consistent, pragmatic style.

## Dead Code

Remove unused code immediately. Do not leave commented-out code, unused variables, unreferenced functions, or stale imports. If something is no longer needed, delete it.

## DRY — Pragmatically

Avoid duplicating logic, but do not create abstractions preemptively. Extract shared logic when the same pattern appears in at least two or three places and the extraction genuinely simplifies the code. A little repetition is preferable to an over-engineered abstraction.

## No Unnecessary Abstractions

Do not hide an implementation behind an interface when there is only one implementation. Use concrete types directly. Introduce interfaces only when you have multiple implementations or are writing code that genuinely needs to be tested with a mock/stub.

## Naming Over Comments

Prefer clear, descriptive names for variables, functions, and types over short names with accompanying comments. A longer, self-explanatory name is better than an abbreviated one that requires explanation. Reserve comments for non-obvious decisions or external constraints — not for describing what the code does.


## Project Structure

The application is one Go module with three binaries under `cmd/`, all built on
one shared package tree under `internal/`.

```
cmd/web/          server-rendered htmx + Alpine app
cmd/api/          JSON transport
cmd/mcp/          MCP transport (streamable HTTP)

tools/copyassets.mjs  build step: static assets and npm libs into tmp/

internal/service/          domain types and operations
internal/web/              everything cmd/web needs
  features/                one subdirectory per user-facing feature
  components/              UI building blocks with no HTTP surface
  static/                  stylesheets, images, fonts
internal/api/              JSON handlers
internal/mcpserver/        MCP tool definitions
internal/infrastructure/   platform and runtime concerns (config, logging,
                           health, metrics, compression, static files)
```

`cmd/` holds only what gets deployed, so `./cmd/...` is the whole shipped
surface. Programs that exist to build the project live in `tools/` and are never
part of an image.

### The rule that matters

**`internal/service` imports nothing else under `internal/`.** It holds the
domain types and operations; `web`, `api` and `mcpserver` all point inward at
it. Each transport is thin — it translates a request into a service call and a
service result into its own idiom. If `service` ever needs to import a
transport or a store, the dependency has gone the wrong way round.

The web app binds the service directly. It does not call its own API: the Go
interface is already the contract, and routing page renders through HTTP would
add a network hop and a second failure mode for nothing.

Domain errors are UI states, so they are typed values (`service.ErrEmptyName`)
that each transport translates — a 400 in `api`, a tool-level error in
`mcpserver`, a message on the page in `web`. A refusal must never become a 500.

### Placement

- Has a route and renders HTML? A feature, in `internal/web/features/`.
- Has a route and returns JSON? A handler in `internal/api/`.
- An LLM-callable tool? `internal/mcpserver/`.
- A UI building block with no HTTP surface? `internal/web/components/`.
- Domain type or operation, with no transport in it? `internal/service/`.
- Platform or runtime concern, with no feature or UI logic? `internal/infrastructure/`.

Do not add new top-level directories under `internal/` without good reason.

Feature-internal components (used only within one feature) live in the feature
directory and are unexported. Components used across features live in
`internal/web/components/`.

### Showcase

Every component in `internal/web/components/` has a showcase page rendering it in isolation from any feature. The demo lives *in the component's own directory* as `showcase.go` (plus any templates it needs), exporting `var Showcase = demo.Page{...}`. A component and its demo therefore change together, in the same diff.

`internal/web/features/showcase/` owns only the routes, the index, and the page chrome; it collects the exported `Showcase` values into its running order. Demos for stylesheets rather than components — typography, colours, spacing, buttons, menu, meter — have no package to sit alongside, so they stay in the feature.

When you add or change a shared component, add or update its showcase page in the same change. It is where the component gets designed and reviewed.

## Language-Specific Guidelines

- Go: see [agents/golang.md](agents/golang.md)
- HTML: see [agents/html.md](agents/html.md)
- JavaScript: see [agents/js.md](agents/js.md)

## Testing

Add a test when there is a genuine reason: the function has multiple edge cases that are non-obvious, the output is hard to verify through normal use, or a bug has been fixed and regression coverage is valuable. Do not test functions simply to confirm they work — if the behaviour is obvious and a manual run through the app would surface any breakage, a test adds noise without value. When tests are warranted, use table-driven tests for functions with multiple input/output cases.

## Verifying Changes

Do not launch the web server (`go run ./cmd/web`) to verify a change. A dev instance is normally already running on the default port; a second `go run` will either fail to bind or fight over it. Verify with `go build ./...` / `go vet ./...` (or `make test`), and ask the user to check the running instance for anything that needs visual or browser confirmation.

`cmd/api` and `cmd/mcp` have no dev instance, so they can be run directly to check a change — they default to ports 8081 and 8082.

## Keep It Simple

Do not over-engineer. Solve the problem at hand. Do not add configuration, flags, or extension points for requirements that do not yet exist. The right amount of complexity is the minimum needed for the current task.
