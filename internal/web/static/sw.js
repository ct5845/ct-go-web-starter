const CACHE = "ct-go-web-starter-cache-v1";
const OFFLINE_URL = "/static/offline.html";

const PRECACHE = ["/static/style.css", OFFLINE_URL];

// Install: cache static assets + offline page
self.addEventListener("install", (e) => {
  e.waitUntil(
    caches
      .open(CACHE)
      .then((cache) => cache.addAll(PRECACHE))
      .then(() => self.skipWaiting()),
  );
});

// Activate: delete old caches
self.addEventListener("activate", (e) => {
  e.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(
          keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)),
        ),
      )
      .then(() => self.clients.claim()),
  );
});

// Fetch: static assets from cache, everything else from network
self.addEventListener("fetch", (e) => {
  const url = new URL(e.request.url);

  // Static assets: cache-first
  if (url.pathname.startsWith("/static/")) {
    e.respondWith(
      caches.match(e.request).then((cached) => cached || fetch(e.request)),
    );
    return;
  }

  // Pages: network-only, offline fallback
  e.respondWith(fetch(e.request).catch(() => caches.match(OFFLINE_URL)));
});
