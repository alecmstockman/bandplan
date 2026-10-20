console.log("notification-websocket.js loaded");

let reconnectAttempts = 0;
let reconnectTimer = null;
let shouldReconnect = true;

function handleIncomingNotification(event) {
	console.log("handleIncomingMessage")
	let message;

	try {
		message = JSON.parse(event.data);
	} catch (error) {
		console.error("Unable to parse WebSocket notification:", error);
		return;
	}

	if (message.type !== "notification_message") {
		return;
	}

	const notificationsElement = document.getElementById("notifications");

	if (!notificationsElement) {
		return;
	}

	if (message.chat_id !== notificationsElement.dataset.chatId) {
		return;
	}

	const currentUserId = notificationsElement.dataset.currentUserId;

	if (message.user_id === currentUserId) {
		appendOwnMessage(notificationsElement, message);
	} else {
		appendOtherMessage(notificationsElement, message);
	}

}


function setNotificationFormConnected(isConnected) {
	console.log("setNotificationFormConnected")

    const input = document.getElementById("message-input");
    const button = document.querySelector(
        "#message-form button[type='submit']",
    );

    if (input) {
        input.disabled = !isConnected;
        input.placeholder = isConnected
            ?"Type a message"
            : "Reconnecting...";
    }

    if (button) {
        button.disabled = !isConnected;
    }
}

function connectNotificationWebSocket() {
	console.log("connectNotificationWebSocket")
    setNotificationFormConnected(false);

	const protocol =
		window.location.protocol === "https:" ? "wss:" : "ws:";

	const socketURL =
		`${protocol}//${window.location.host}/ws/chat`;

	const socket = new WebSocket(socketURL);

	socket.addEventListener("open", () => {
		console.log("Notification WebSocket connected");

        reconnectAttempts = 0;

        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }

        setNotificationFormConnected(true);
	});

	socket.addEventListener("message", handleIncomingMessage);

	socket.addEventListener("close", (event) => {
		console.log("Chat WebSocket disconnected:", event.code);

        setChatFormConnected(false);

        if (shouldReconnect) {
            scheduleReconnect();
        }
	});

	socket.addEventListener("error", (error) => {
		console.error("Chat WebSocket error:", error);
	});

	return socket;
}
