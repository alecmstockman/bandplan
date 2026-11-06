let notificationReconnectAttempts = 0;
let notificationReconnectTimer = null;
let notificationShouldReconnect = true;

function displayNotification(notificationElement) {
    let icon = document.getElementById("notification-icon");

    if (!icon) {
        icon = document.createElement("div");
        icon.id = "notification-icon";
        icon.textContent = "0";
        notificationElement.prepend(icon);
    }

    const currentCount = Number.parseInt(icon.textContent, 10) || 0;
    icon.textContent = String(currentCount + 1);
}

function handleIncomingNotification(event) {
    let notification;

    try {
        notification = JSON.parse(event.data);
    } catch (error) {
        console.error("Unable to parse WebSocket notification:", error);
        return;
    }

    if (notification.type !== "notification_message") {
        return;
    }

    const notificationElement = document.getElementById("notification-box");
    if (notificationElement) {
        displayNotification(notificationElement);
    }
}

function scheduleNotificationReconnect() {
    if (notificationReconnectTimer) {
        return;
    }

    const delay = Math.min(
        1000 * Math.pow(2, notificationReconnectAttempts),
        30000,
    );
    notificationReconnectAttempts++;

    notificationReconnectTimer = window.setTimeout(() => {
        notificationReconnectTimer = null;
        window.notificationSocket = connectNotificationWebSocket();
    }, delay);
}

function connectNotificationWebSocket() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const socketURL = `${protocol}//${window.location.host}/ws/notifications`;
    const socket = new WebSocket(socketURL);

    socket.addEventListener("open", () => {
        notificationReconnectAttempts = 0;

        if (notificationReconnectTimer) {
            window.clearTimeout(notificationReconnectTimer);
            notificationReconnectTimer = null;
        }
    });

    socket.addEventListener("message", handleIncomingNotification);

    socket.addEventListener("close", () => {
        if (notificationShouldReconnect) {
            scheduleNotificationReconnect();
        }
    });

    socket.addEventListener("error", (error) => {
        console.error("Notification WebSocket error:", error);
    });

    return socket;
}

if (document.getElementById("notification-box")) {
    window.notificationSocket = connectNotificationWebSocket();
}

window.addEventListener("beforeunload", () => {
    notificationShouldReconnect = false;

    if (notificationReconnectTimer) {
        window.clearTimeout(notificationReconnectTimer);
        notificationReconnectTimer = null;
    }

    if (window.notificationSocket) {
        window.notificationSocket.close(1000, "Page unloading");
    }
});
