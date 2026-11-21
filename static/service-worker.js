const CACHE_NAME = "bandplan-v2";

const STATIC_ASSETS = [
  "/static/manifest.webmanifest",
  "/static/images/icons/icon-192.png",
  "/static/images/icons/icon-512.png",
  "/static/images/icons/apple-touch-icon.png"
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => cache.addAll(STATIC_ASSETS))
      .then(() => self.skipWaiting())
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((cacheNames) => Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      ))
      .then(() => self.clients.claim())
  );
});

self.addEventListener("push", (event) => {
  let data = {};

  try {
    data = event.data?.json() ?? {};
  } catch {
    data.body = event.data?.text() ?? "You have a notification.";
  }

  event.waitUntil(
    self.registration.showNotification(data.title ?? "BandPlan", {
      body: data.body ?? "You have a notification.",
      icon: "/static/images/icons/icon-192.png",
      data: {
        url: data.url ?? "/notifications",
      },
    })
  );
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();

  const url = event.notification.data?.url ?? "/notifications";

  event.waitUntil(self.clients.openWindow(url));
});
