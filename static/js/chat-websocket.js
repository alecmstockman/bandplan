console.log("chat-websocket.js loaded");

function appendOwnMessage(messagesElement, message) {
	const listItem = document.createElement("li");
	listItem.className = "message-own";

	const body = document.createElement("div");
	body.className = "message-body";
	body.append(document.createTextNode(message.body));

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
	body.append(document.createTextNode(message.body));

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
	const protocol =
		window.location.protocol === "https:" ? "wss:" : "ws:";

	const socketURL =
		`${protocol}//${window.location.host}/ws/chat`;

	console.log("WebSocket URL:", socketURL);

	const socket = new WebSocket(socketURL);

	socket.addEventListener("open", () => {
		console.log("Chat WebSocket connected");
	});

	socket.addEventListener("message", handleIncomingMessage);

	socket.addEventListener("close", (event) => {
		console.log("Chat WebSocket disconnected:", event.code);
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

if (
	!window.chatSocket ||
	window.chatSocket.readyState === WebSocket.CLOSED
) {
	window.chatSocket = connectChatWebSocket();
}

configureMessageForm();