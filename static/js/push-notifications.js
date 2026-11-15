const enablePush = document.querySelector("#enable-push");
const pushStatus = document.querySelector("#push-status");

enablePush.addEventListener("click", async () => {
    try {
        const registration = await navigator.serviceWorker.ready;
        let subscription = await registration.pushManager.getSubscription();

        if (!subscription) {
            subscription = await registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: enablePush.dataset.vapidPublicKey,
            });
        }

        pushStatus.textContent = "Browser subscription created.";
        console.log("Push providers:", new URL(subscription.endpoint).origin);
    } catch (error) {
        pushStatus.textContent = "Could not create a push subscription.";
        console.error(error);
    }
});