# Go Guidelines

## Idioms and Standard Library

Write idiomatic Go. Follow standard Go conventions for naming, error handling, and package organisation. Before reaching for a third-party library, ask whether the standard library (`net/http`, `encoding/json`, `html/template`, etc.) covers the need. Add a dependency only when it provides significant, non-trivial value over stdlib.

File names use lowercase with no hyphens — concatenate words directly (e.g. `welcomecard.go`, not `welcome-card.go`). Underscores are reserved for test files (`foo_test.go`) and platform-specific build files (`foo_windows.go`).

Acronyms stay all-caps: `userID` not `userId`, `parseURL` not `parseUrl`, `serveHTTP` not `serveHttp`.

Receiver names are short, derived from the type: `func (c *component)` not `func (comp *component)` or `func (this *component)`.

No `Get` prefix on getters: `user.Name()` not `user.GetName()`.

Always name the error variable `err`. Redeclare with `:=` in new scopes rather than inventing `err2` or `e`.

Sentinel errors use `var`: `var ErrNotFound = errors.New("not found")`.

Only use `New` as a constructor prefix when initialisation is non-trivial. A plain struct literal or function with a clear name is preferable otherwise.

Do not create `utils`, `helpers`, `types`, or `models` packages. Name packages by what they provide. Keep types next to the code that owns them.

Before writing new helper code, check whether a package already exists in `src/` that covers the need. Reuse it rather than duplicating locally.

## Prefer Functions Over Methods

Prefer package-level functions over methods on structs where there is no meaningful state to encapsulate. A struct with no real state that exists only to hang methods off is an unnecessary indirection — use a plain function instead. Use structs and methods when the type genuinely owns state that needs to travel with behaviour.

## Error Handling

Return errors to the caller; do not swallow them silently. Do not log an error and then also return it — pick one. Wrap errors with context using `fmt.Errorf("doing X: %w", err)` so call sites have enough information. At HTTP boundaries, translate errors into appropriate status codes and log once.

## Panics

Panic only for unrecoverable programmer errors at initialisation time (e.g. a template that fails to parse on startup). Never panic in request-handling code — return an error instead.

## Logging

Log meaningful events at appropriate levels. Avoid noisy debug logs that restate what the function name already says. Log at `slog.Info` for significant lifecycle events, `slog.Warn` for unexpected-but-recoverable situations, and `slog.Error` when something fails. Include relevant structured fields, not prose descriptions of the code path.

## Routing

All routes are registered in `src/app.go` using `http.NewServeMux()` from the standard library. Do not introduce a third-party router. Each feature exposes a single handler function (e.g. `home.Handler`) which is registered directly on the mux. Do not create a separate router file or a route registration abstraction — just add the `mux.HandleFunc` call in `App()`.


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


When a feature handler has meaningful page assembly work (rendering subcomponents, preparing data), split it into two files:

- `handler.go` — HTTP only: validate the request, call `renderPage()`, write the response or error
- `page.go` — UI assembly: embed templates, render subcomponents, compose and return the full page HTML

Only split when there is real assembly work. A trivial handler with no subcomponents does not need a separate `page.go`.

See `src/features/home/` for a working example of this pattern.

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
// memberBottomTabs is computed once and reused across handlers.
var memberBottomTabs = bottomtabs.MustRender(bottomtabs.Options{...})
```
