const CACHE_NAME = "bandplan-v1";

const STATIC_ASSETS = [
  "/static/manifest.json",
  "/static/images/icons/icon-192.png",
  "/static/images/icons/icon-512.png",
  "/static/images/icons/apple-touch-icon.png"
];

self.addEventListener("install", event => {
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => {
      return cache.addAll([
        "/",
        "/static/manifest.json",
        "/static/images/icons/icon-192.png",
        "/static/images/icons/icon-512.png",
        "/static/images/icons/apple-touch-icon.png"
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames
          .filter(name => name !== CACHE_NAME)
          .map(name => caches.delete(name))
      );
    })
  );
});