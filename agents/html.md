# HTML Guidelines

## Accessibility

Accessibility is a first-class requirement, not an afterthought. Apply it from the start, not as a pass at the end.

**Semantics first.** Use the correct HTML element for the job — `<button>` for actions, `<a>` for navigation, `<nav>`, `<main>`, `<section>`, `<article>` for structure. A `<div>` with an `@click` is never the right answer when a `<button>` exists.

**ARIA only when native semantics fall short.** Don't add `role="button"` to a `<button>` — it's redundant. Do use ARIA when you build something the browser has no native element for. Key attributes to consider:
- `aria-label` / `aria-labelledby` — when visible text doesn't describe the element (e.g. an icon-only button)
- `aria-hidden="true"` — on decorative elements (icons, spinners) that should be invisible to screen readers
- `aria-live="polite"` — on regions that update dynamically (status messages, streaming responses)
- `aria-expanded`, `aria-selected`, `aria-disabled` — for interactive state that isn't conveyed by the element type alone

**Showing and hiding.** Prefer `x-show` (or CSS `visibility`/`display`) over inserting/removing elements from the DOM for dynamic content. A persistent element with `aria-live` is more reliably announced by screen readers than one that appears suddenly. When you do insert content dynamically, make sure focus and announcement behaviour is intentional.

**Focus management.** When an interaction moves the user to a new context (opens a modal, navigates to a new view), move focus explicitly. When the context closes, return focus to the trigger. Never leave focus stranded on a removed element.

## SEO & GEO

Make sure where possible we consider SEO and GEO.

## Plain Prototyping

When prototyping and generating HTML don't create random classes and styles unless it specifically is to fit the brief. Let it be semantically correct, that's all.

## Prefer Native HTML

Always reach for the simplest, most native solution first. Before adding JavaScript or a framework, ask whether plain HTML achieves the same result. Examples: `autofocus` over `element.focus()`, `required` over JS validation, `<details>` over a JS toggle, `<dialog>` over a custom modal. Native HTML is more accessible, more resilient, and requires no maintenance.