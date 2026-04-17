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
