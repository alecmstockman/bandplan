const enablePush = document.querySelector("#enable-push");
const pushStatus = document.querySelector("#push-status");

enablePush?.addEventListener("click", async () => {
    try {
        const registration = await navigator.serviceWorker.ready;
        let subscription = await registration.pushManager.getSubscription();

        if (!subscription) {
            subscription = await registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: enablePush.dataset.vapidPublicKey,
            });
        }

        const response = await fetch("/push/subscription", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-CSRF-Token": enablePush.dataset.csrfToken,
            },
            body: JSON.stringify(subscription.toJSON()),
        });

        if (response.status !== 204) {
            throw new Error(`Could not save push subscription (${response.status})`);
        }

        pushStatus.textContent = "Push subscription saved.";
    } catch (error) {
        pushStatus.textContent = "Could not enable push notifications.";
        console.error(error);
    }
});
