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

## Routing

All routes are registered in `src/app.go` using `http.NewServeMux()` from the standard library. Do not introduce a third-party router. Each feature exposes a single handler function (e.g. `home.Handler`) which is registered directly on the mux. Do not create a separate router file or a route registration abstraction — just add the `mux.HandleFunc` call in `App()`.

## Project Structure

`src/` has three top-level directories. Keep it that way — do not add new top-level directories without good reason.

- `features/` — one subdirectory per user-facing feature (e.g. `features/home/`)
- `shared/` — UI components and templates reused across features (`shared/component/`, `shared/components/`, `shared/templates/`)
- `infrastructure/` — platform and runtime concerns with no feature or UI logic (`infrastructure/fs/`, `infrastructure/fileserver/`, `infrastructure/config/`)

If something is only used by one feature, it lives in that feature's directory. If it's genuinely reused across features, it goes in `shared/`. If it's a platform/runtime concern (file I/O, HTTP serving, config), it goes in `infrastructure/`.

## Components and Templates

Each component is a `.go` file + `.html` file pair, optionally with a `.js` file when using `component.WithJS`. Use `//go:embed` to embed the HTML at compile time.

All data preparation happens in Go. Do not use template functions for logic. If a template needs data, compute it in Go and pass it as a named prop. Templates are for rendering only.

Subcomponents are rendered in Go first, then passed to the parent template as `template.HTML` props. Never render a component from inside another template — compose in Go, not in HTML.

```
// render the subcomponent in Go
welcomeCardHTML, err := renderWelcomeCard("Title", "Description")

// pass it as a prop to the parent
homeTmpl.Render("WelcomeCardHTML", welcomeCardHTML)
```

```html
<!-- parent template receives it as a plain prop -->
{{ .WelcomeCardHTML }}
```

Feature-internal components (used only within one feature) live in the feature directory and are unexported. Shared components used across features live in `src/shared/components/`.

When a feature handler has meaningful page assembly work (rendering subcomponents, preparing data), split it into two files:

- `handler.go` — HTTP only: validate the request, call `renderPage()`, write the response or error
- `page.go` — UI assembly: embed templates, render subcomponents, compose and return the full page HTML

Only split when there is real assembly work. A trivial handler with no subcomponents does not need a separate `page.go`.

See `src/features/home/` for a working example of this pattern.

## Go Idioms and Standard Library

Write idiomatic Go. Follow standard Go conventions for naming, error handling, and package organisation. Before reaching for a third-party library, ask whether the standard library (`net/http`, `encoding/json`, `html/template`, etc.) covers the need. Add a dependency only when it provides significant, non-trivial value over stdlib.

File names use lowercase with no hyphens — concatenate words directly (e.g. `welcomecard.go`, not `welcome-card.go`). Underscores are reserved for test files (`foo_test.go`) and platform-specific build files (`foo_windows.go`).

Acronyms stay all-caps: `userID` not `userId`, `parseURL` not `parseUrl`, `serveHTTP` not `serveHttp`.

Receiver names are short, derived from the type: `func (c *component)` not `func (comp *component)` or `func (this *component)`.

No `Get` prefix on getters: `user.Name()` not `user.GetName()`.

Always name the error variable `err`. Redeclare with `:=` in new scopes rather than inventing `err2` or `e`.

Sentinel errors use `var`: `var ErrNotFound = errors.New("not found")`.

Only use `New` as a constructor prefix when initialisation is non-trivial. A plain struct literal or function with a clear name is preferable otherwise.

Do not create `utils`, `helpers`, `types`, or `models` packages. Name packages by what they provide. Keep types next to the code that owns them. For example, file system helpers live in `src/infrastructure/fs/`, not `src/utils/`.

`src/shared/` is an organisational directory, not a package — the packages inside it (`component`, `components`, `templates`) are named by what they provide. This is fine. The word "shared" describes location, not purpose.

Before writing new helper code, check whether a package already exists in `src/` that covers the need. Reuse it rather than duplicating locally.

## Prefer Functions Over Methods

Prefer package-level functions over methods on structs where there is no meaningful state to encapsulate. A struct with no real state that exists only to hang methods off is an unnecessary indirection — use a plain function instead. Use structs and methods when the type genuinely owns state that needs to travel with behaviour.

## Error Handling

Return errors to the caller; do not swallow them silently. Do not log an error and then also return it — pick one. Wrap errors with context using `fmt.Errorf("doing X: %w", err)` so call sites have enough information. At HTTP boundaries, translate errors into appropriate status codes and log once.

## Panics

Panic only for unrecoverable programmer errors at initialisation time (e.g. a template that fails to parse on startup). Never panic in request-handling code — return an error instead.

## Logging

Log meaningful events at appropriate levels. Avoid noisy debug logs that restate what the function name already says (e.g. "Executing component template", "Component created successfully"). Log at `slog.Info` for significant lifecycle events, `slog.Warn` for unexpected-but-recoverable situations, and `slog.Error` when something fails. Include relevant structured fields, not prose descriptions of the code path.

## Testing

Do not write tests by default. Add a test when there is a genuine reason: the function has multiple edge cases that are non-obvious, the output is hard to verify through normal use, or a bug has been fixed and regression coverage is valuable. Do not test functions simply to confirm they work — if the behaviour is obvious and a manual run through the app would surface any breakage, a test adds noise without value. When tests are warranted, use table-driven tests for functions with multiple input/output cases.

## Keep It Simple

Do not over-engineer. Solve the problem at hand. Do not add configuration, flags, or extension points for requirements that do not yet exist. The right amount of complexity is the minimum needed for the current task.
