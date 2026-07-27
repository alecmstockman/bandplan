console.log("chat-websocket.js loaded");

let reconnectAttempts = 0;
let reconnectTimer = null;
let shouldReconnect = true;

function setChatFormConnected(isConnected) {
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

function appendOwnMessage(messagesElement, message) {
	const listItem = document.createElement("li");
	listItem.className = "message-own";

	const body = document.createElement("div");
	body.className = "message-body";
	// body.append(document.createTextNode(message.body));
    appendTextWithLinks(body, message.body);

	const footer = document.createElement("div");
	footer.className = "message-body-footer";
	footer.textContent = message.display_time;

	body.appendChild(footer);
	listItem.appendChild(body);
	messagesElement.appendChild(listItem);
}

function appendOtherMessage(messagesElement, message) {
	const listItem = document.createElement("li");
	listItem.className = "test-message-other";

	const pictureBox = document.createElement("div");
	pictureBox.className = "test-message-sender-pic-box";

	const picture = document.createElement("img");
	picture.className = "test-message-sender-pic";
	picture.src = message.profile_image_path || "";
	picture.alt = "";

	pictureBox.appendChild(picture);

	const content = document.createElement("div");
	content.className = "test-message-content";

	const header = document.createElement("div");
	header.className = "message-header";
	header.textContent = message.user_name;

	const body = document.createElement("div");
	body.className = "message-body-other";
	// body.append(document.createTextNode(message.body));
    appendTextWithLinks(body, message.body);

	const footer = document.createElement("div");
	footer.className = "message-body-footer";
	footer.textContent = message.display_time;

	body.appendChild(footer);
	content.appendChild(header);
	content.appendChild(body);

	listItem.appendChild(pictureBox);
	listItem.appendChild(content);

	messagesElement.appendChild(listItem);
}

function handleIncomingMessage(event) {
	let message;

	try {
		message = JSON.parse(event.data);
	} catch (error) {
		console.error("Unable to parse WebSocket message:", error);
		return;
	}

	if (message.type !== "chat_message") {
		return;
	}

	const messagesElement = document.getElementById("messages");

	if (!messagesElement) {
		return;
	}

	const currentUserID = messagesElement.dataset.currentUserId;

	if (message.user_id === currentUserID) {
		appendOwnMessage(messagesElement, message);
	} else {
		appendOtherMessage(messagesElement, message);
	}

	messagesElement.scrollTop = messagesElement.scrollHeight;
}

function connectChatWebSocket() {
    setChatFormConnected(false);

	const protocol =
		window.location.protocol === "https:" ? "wss:" : "ws:";

	const socketURL =
		`${protocol}//${window.location.host}/ws/chat`;

	console.log("WebSocket URL:", socketURL);

	const socket = new WebSocket(socketURL);

	socket.addEventListener("open", () => {
		console.log("Chat WebSocket connected");

        reconnectAttempts = 0;

        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null
        }

        setChatFormConnected(true);
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

function configureMessageForm() {
	const form = document.getElementById("message-form");
	const input = document.getElementById("message-input");

	if (!form || !input) {
		return;
	}

	form.addEventListener("submit", (event) => {
		event.preventDefault();

		const body = input.value.trim();

		if (!body) {
			return;
		}

		if (
			!window.chatSocket ||
			window.chatSocket.readyState !== WebSocket.OPEN
		) {
			console.error("Chat WebSocket is not connected");
			return;
		}

        console.log("Sending WebSocket message:", body);
        console.log("Socket state:", window.chatSocket.readyState);

		window.chatSocket.send(
			JSON.stringify({
				type: "chat_message",
				body: body,
			}),
		);

		form.reset();
		input.focus();
	});
}

function scheduleReconnect() {
    if (reconnectTimer) {
        return;
    }

    const delay = Math.min(
        1000 * Math.pow(2, reconnectAttempts),
        30000,
    );

    reconnectAttempts++;

    console.log(`Reconnecting in ${delay}ms`);

    reconnectTimer = setTimeout(() => {
        reconnectTimer = null;

        window.chatSocket = connectChatWebSocket();
    }, delay);
}

if (
	!window.chatSocket ||
	window.chatSocket.readyState === WebSocket.CLOSED
) {
	window.chatSocket = connectChatWebSocket();
}




window.addEventListener("beforeunload", () => {
	shouldReconnect = false;

	if (reconnectTimer) {
		clearTimeout(reconnectTimer);
		reconnectTimer = null;
	}

	if (window.chatSocket) {
		window.chatSocket.close(1000, "Page unloading");
	}
});

function scrollMessagesToBottom() {
	const messagesElement = document.getElementById("messages");

	if (!messagesElement) {
		return;
	}

	messagesElement.scrollTop = messagesElement.scrollHeight;
}

function appendTextWithLinks(container, text) {
	const urlPattern = /(https?:\/\/[^\s]+)/g;
	const parts = text.split(urlPattern);

	for (const part of parts) {
		if (urlPattern.test(part)) {
			const link = document.createElement("a");

			link.href = part;
			link.textContent = part;
			link.target = "_blank";
			link.rel = "noopener noreferrer";

			container.appendChild(link);
		} else {
			container.appendChild(document.createTextNode(part));
		}

		urlPattern.lastIndex = 0;
	}
}

function linkifyExistingMessages() {
	const messageTextElements =
		document.querySelectorAll(".message-text");

	for (const element of messageTextElements) {
		const text = element.textContent;

		element.textContent = "";
		appendTextWithLinks(element, text);
	}
}



configureMessageForm();
linkifyExistingMessages();
scrollMessagesToBottom();
