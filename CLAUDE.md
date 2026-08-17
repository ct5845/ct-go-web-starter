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

`src/` has three top-level directories. Keep it that way — do not add new top-level directories without good reason.

- `features/` — one subdirectory per user-facing feature (e.g. `features/home/`); features have an HTTP surface (handler + routes)
- `components/` — UI building blocks with no HTTP surface (`components/component/`, `components/page/`, `components/icon/`, etc.)
- `infrastructure/` — platform and runtime concerns with no feature or UI logic (`infrastructure/fileserver/`, `infrastructure/config/`, `infrastructure/compression/`)

The rule for placement is simple: if it has a route, it's a feature. If it's a UI building block with no HTTP surface, it's a component. If it's a platform/runtime concern, it's infrastructure.

Feature-internal components (used only within one feature) live in the feature directory and are unexported. Components used across features live in `src/components/`.

Every component in `src/components/` has a showcase page rendering it in isolation from any feature. The demo lives *in the component's own directory* as `showcase.go` (plus any templates it needs), exporting `var Showcase = demo.Page{...}`. A component and its demo therefore change together, in the same diff.

`src/features/showcase/` owns only the routes, the index, and the page chrome; it collects the exported `Showcase` values into its running order. Demos for stylesheets rather than components — typography, colours, spacing, buttons, menu, meter — have no package to sit alongside, so they stay in the feature.

When you add or change a shared component, add or update its showcase page in the same change. It is where the component gets designed and reviewed.

## Language-Specific Guidelines

- Go: see [agents/golang.md](agents/golang.md)
- HTML: see [agents/html.md](agents/html.md)
- JavaScript: see [agents/js.md](agents/js.md)

## Testing

Do not write tests by default. Add a test when there is a genuine reason: the function has multiple edge cases that are non-obvious, the output is hard to verify through normal use, or a bug has been fixed and regression coverage is valuable. Do not test functions simply to confirm they work — if the behaviour is obvious and a manual run through the app would surface any breakage, a test adds noise without value. When tests are warranted, use table-driven tests for functions with multiple input/output cases.

## Verifying Changes

Do not launch the web server (`go run ./cmd/web`) to verify a change. A dev instance is normally already running on the default port; a second `go run` will either fail to bind or fight over it. Verify with `go build ./...` / `go vet ./...`, and ask the user to check the running instance for anything that needs visual or browser confirmation.

## Keep It Simple

Do not over-engineer. Solve the problem at hand. Do not add configuration, flags, or extension points for requirements that do not yet exist. The right amount of complexity is the minimum needed for the current task.
