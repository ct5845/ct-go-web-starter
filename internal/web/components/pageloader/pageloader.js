const bar = document.getElementById("pageloader");

let timer = null;
let progress = 0;

function isPlainClick(event) {
  return event.button === 0 && !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey;
}

function isNavigableLink(link) {
  if (!link || !link.href) return false;
  if (link.target && link.target !== "_self") return false;
  if (link.hasAttribute("download")) return false;
  if (link.origin !== window.location.origin) return false;
  if (link.hash && link.pathname === window.location.pathname && link.search === window.location.search) return false;
  return true;
}

function start() {
  clearInterval(timer);
  progress = 0;
  bar.style.transition = "none";
  bar.style.width = "0%";
  bar.style.opacity = "1";
  // Force layout so the next transition (triggered by the interval below)
  // animates from 0 instead of jumping straight to the first tick's width.
  bar.offsetWidth;
  bar.style.transition = "width 200ms ease-out, opacity 150ms ease-in";

  timer = setInterval(() => {
    // Approach 90% asymptotically so the bar never claims completion before
    // the real navigation or request actually finishes.
    progress += (90 - progress) * 0.1;
    bar.style.width = progress + "%";
  }, 200);
}

function finish() {
  clearInterval(timer);
  bar.style.width = "100%";
  setTimeout(() => {
    bar.style.opacity = "0";
  }, 150);
}

document.addEventListener("click", (event) => {
  if (!isPlainClick(event)) return;
  const link = event.target.closest("a");
  if (!link || link.hasAttribute("hx-boost") || link.closest("[hx-boost]")) return;
  if (!isNavigableLink(link)) return;
  start();
});

document.addEventListener("submit", (event) => {
  const form = event.target;
  if (form.hasAttribute("hx-post") || form.hasAttribute("hx-get") || form.closest("[hx-boost]")) return;
  start();
});

document.body.addEventListener("htmx:beforeRequest", (event) => {
  if (!event.detail.boosted) return;
  start();
});
document.body.addEventListener("htmx:afterRequest", (event) => {
  if (!event.detail.boosted) return;
  finish();
});

// bfcache restores (back/forward) don't re-run the IIFE's module-load state
// but do fire pageshow, so reset explicitly rather than leaving a stuck bar.
window.addEventListener("pageshow", () => {
  clearInterval(timer);
  bar.style.transition = "none";
  bar.style.width = "0%";
  bar.style.opacity = "0";
});
