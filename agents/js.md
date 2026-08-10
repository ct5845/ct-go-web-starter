# JavaScript Guidelines

## Prefer the Simplest Tool

Reach for the most native, least complex solution first. Work up the stack only when the simpler option genuinely cannot do the job:

1. **Plain HTML** — if a native element or attribute solves it, use it (`autofocus`, `required`, `<details>`, `<dialog>`)
2. **CSS** — if it's visual or state-driven presentation, use CSS (`:hover`, `:focus`, `@keyframes`, `transition`)
3. **Vanilla JS** — for DOM interaction, events, and fetch that don't need reactivity
4. **Alpine.js** — for reactive UI state that would be awkward to manage manually
5. **Anything heavier** — only with strong justification; we don't have a build step by default

Don't jump to Alpine (or any reactive layer) for something a single HTML attribute or a CSS rule handles. The cost of the wrong abstraction is complexity that compounds.

## Alpine.js

Use Alpine for reactive state — things that change at runtime and need to update the UI. Keep the Alpine component focused on state and behaviour; derive display values via getters rather than scattering logic across handlers.

Prefer a `get computedProp()` getter over calling a method multiple times in the template, or over imperatively updating a prop in every handler that could affect it.

For DOM manipulation that bypasses Alpine's reactivity (e.g. streaming text into a node), do it directly via the DOM — don't force it through `x-text` or `x-html` when direct `textContent` mutation is simpler and more performant.

## Event Cleanup

Remove event listeners when they are no longer needed. For one-time events use `{ once: true }`. For EventSource and WebSocket: close the connection and null the reference in a single `stop` method so cleanup is always in one place.
