# Go Guidelines

## Idioms and Standard Library

Write idiomatic Go. Follow standard Go conventions for naming, error handling, and package organisation. Before reaching for a third-party library, ask whether the standard library (`net/http`, `encoding/json`, `html/template`, etc.) covers the need. Modules maintained by the Go team (`golang.org/x/...`) count as stdlib for this purpose — use them freely. Add any other dependency only when it provides significant, non-trivial value over stdlib.

File names use lowercase with no hyphens — concatenate words directly (e.g. `welcomecard.go`, not `welcome-card.go`). Underscores are reserved for test files (`foo_test.go`) and platform-specific build files (`foo_windows.go`).

Acronyms stay all-caps: `userID` not `userId`, `parseURL` not `parseUrl`, `serveHTTP` not `serveHttp`.

Receiver names are short, derived from the type: `func (c *component)` not `func (comp *component)` or `func (this *component)`.

No `Get` prefix on getters: `user.Name()` not `user.GetName()`.

Always name the error variable `err`. Redeclare with `:=` in new scopes rather than inventing `err2` or `e`.

Sentinel errors use `var`: `var ErrNotFound = errors.New("not found")`.

Only use `New` as a constructor prefix when initialisation is non-trivial. A plain struct literal or function with a clear name is preferable otherwise.

Do not create `utils`, `helpers`, `types`, or `models` packages. Name packages by what they provide. Keep types next to the code that owns them.

Before writing new helper code, check whether a package already exists in `src/` that covers the need. Reuse it rather than duplicating locally.

## Comments

Comment the *why*, never the *what*. A comment that restates what the code plainly does (`// Get returns the override for domain`, `// loadGob decodes the file`) is noise — it makes code harder to scan and goes stale. Delete it and let the name and signature carry the meaning.

Keep a comment only when it captures something the code cannot: a constraint (`// field order must match for gob decoding`), a non-obvious rationale (`// clearing primary also clears secondary`), an external quirk, or a deliberate trade-off. If you're tempted to explain what a function does, rename it instead.

This applies to exported identifiers too, including package docs, structs, and their fields. A struct whose field names already say what they hold (`type Options struct { Label string; Hint string }`) needs no doc comment — adding one just restates the names in prose. Write one only when there's a real constraint to capture (an ordering requirement, a unit, a "must match X elsewhere"), and keep it to the one line that says that thing — not a paragraph walking through what every field is. Most types in this codebase should have zero doc comment lines for every line of declaration; if a comment block is longer than the type it describes, that's the signal to cut it.

## Prefer Functions Over Methods

Prefer package-level functions over methods on structs where there is no meaningful state to encapsulate. A struct with no real state that exists only to hang methods off is an unnecessary indirection — use a plain function instead. Use structs and methods when the type genuinely owns state that needs to travel with behaviour.

## Error Handling

Return errors to the caller; do not swallow them silently. Do not log an error and then also return it — pick one. Wrap errors with context using `fmt.Errorf("doing X: %w", err)` so call sites have enough information. At HTTP boundaries, translate errors into appropriate status codes and log once.

## Panics

Panic only for unrecoverable programmer errors at initialisation time (e.g. a template that fails to parse on startup). Never panic in request-handling code — return an error instead.

## Logging

Log meaningful events at appropriate levels. Avoid noisy debug logs that restate what the function name already says. Log at `slog.Info` for significant lifecycle events, `slog.Warn` for unexpected-but-recoverable situations, and `slog.Error` when something fails. Include relevant structured fields, not prose descriptions of the code path.

## Request Tracking

Every operation that does real work on a web request must be wrapped in a `reqlog.Track` span, so the request log shows *where* time goes — not just that a request was slow. Use the deferred form at the top of the scope:

```go
func handleGet(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "home.handleGet", "")()
	// ...
}
```

Add a span to:

- every HTTP handler (keyed `feature.handlerName`);
- every call that crosses an I/O boundary — a database query, an external HTTP call — keyed by the operation. When one logical operation makes several round-trips, track each phase separately so the slow one is visible rather than hidden inside an aggregate.

Pure in-memory work does not need a span. The test is: if this could plausibly be the slow part of a request, it must be trackable in isolation. `Track` is a no-op when the context carries no request, so it is safe to add anywhere.

For streaming responses, where a long wall-clock time is expected rather than a problem, call `reqlog.IgnoreDuration(r.Context())` so the entry is not promoted to Warn.

## Routing

All routes are wired up in `routes()` in `cmd/web/main.go`, using `http.NewServeMux()` from the standard library. Do not introduce a third-party router.

Each feature exposes a single `RegisterRoutes(mux *http.ServeMux)` function that registers its own patterns, and `routes()` calls it:

```go
func routes() *http.ServeMux {
	mux := http.NewServeMux()

	home.RegisterRoutes(mux)
	showcase.RegisterRoutes(mux)
	fileserver.RegisterRoutes(mux, "tmp/static/")

	return mux
}
```

Keeping the patterns inside the feature means a route and its handler stay in one place. Do not create a separate router file or a route registration abstraction beyond this — just add the `RegisterRoutes` call in `routes()`.


## Components and Templates

Each component is a `.go` file + `.html` file pair, optionally with a `.js` file when using `component.WithAlpine` or `component.WithIIFE`. Use `//go:embed` to embed the HTML at compile time.

Pick the wrapper by what the script needs: `WithAlpine` defers the script to the `alpine:init` event, so it can register Alpine stores and data components; `WithIIFE` wraps it in an immediately-invoked function expression, keeping plain DOM scripting out of the global scope.

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


Simple features live in a single file (e.g. `home.go`) that contains the route registration, HTTP handler, and page assembly together. Only split into `handler.go` + `page.go` when there is substantial assembly work — multiple subcomponents, complex data preparation — that would make a single file unwieldy.

See `src/features/home/home.go` for a working example of the single-file pattern.

## Component APIs

Every reusable component takes a typed `Options` struct, not variadic
key/value pairs. Discoverability matters more than terseness — a struct
field shows up in IDE autocomplete and `go doc`; a `keysAndValues` arg
list does not.

```go
// Yes — typed Options.
type Options struct {
    Header     template.HTML
    Content    template.HTML
    BottomTabs template.HTML
}

func Render(options Options) (template.HTML, error) {
    return comp.Render(options)
}

// No — variadic key/value pairs.
func Render(keysAndValues ...any) (template.HTML, error) {
    return comp.Render(keysAndValues...)
}
```

This applies to every component, including page-level ones called only
from a single handler. The handler gains nothing from `keysAndValues`,
and the typed struct documents the contract.

## Composition over branching

A component should not branch internally on a "variant" prop to render
different shapes. Variants are values composed at the call site.

```go
// Yes — composition. Each header is its own value.
homeHeader, _ := header.Render(header.Options{Title: "Home"})
backHeader, _ := header.Render(header.Options{Title: "Sign up", BackURL: "/book"})

// No — internal branching on a Variant prop.
header.Render(header.Options{Variant: "home", Title: "Home"})
header.Render(header.Options{Variant: "back", Title: "Sign up", BackURL: "/book"})
```

If the same composed result appears in many places, store it as a value
and reuse the value — don't push the variation into the component.

```go
// primaryItems is declared once and projected into both the bottom tabs and
// the sidebar, so the two navigations cannot drift apart.
var primaryItems = []primaryItem{...}
```

See `src/features/nav/homenav.go` for this pattern in use.
