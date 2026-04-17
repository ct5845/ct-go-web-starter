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
- `components/` — UI building blocks with no HTTP surface (`components/component/`, `components/header/`, `components/footer/`, `components/page/`)
- `infrastructure/` — platform and runtime concerns with no feature or UI logic (`infrastructure/fs/`, `infrastructure/fileserver/`, `infrastructure/config/`)

The rule for placement is simple: if it has a route, it's a feature. If it's a UI building block with no HTTP surface, it's a component. If it's a platform/runtime concern, it's infrastructure.

Feature-internal components (used only within one feature) live in the feature directory and are unexported. Components used across features live in `src/components/`.

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

## Language-Specific Guidelines

- Go: see [agents/golang.md](agents/golang.md)

## Testing

Do not write tests by default. Add a test when there is a genuine reason: the function has multiple edge cases that are non-obvious, the output is hard to verify through normal use, or a bug has been fixed and regression coverage is valuable. Do not test functions simply to confirm they work — if the behaviour is obvious and a manual run through the app would surface any breakage, a test adds noise without value. When tests are warranted, use table-driven tests for functions with multiple input/output cases.

## Keep It Simple

Do not over-engineer. Solve the problem at hand. Do not add configuration, flags, or extension points for requirements that do not yet exist. The right amount of complexity is the minimum needed for the current task.
