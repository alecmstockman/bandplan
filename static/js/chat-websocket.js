console.log("chat-websocket.js loaded");

let reconnectAttempts = 0;
let reconnectTimer = null;
let shouldReconnect = true;

function setChatFormConnected(isConnected) {
	console.log("setChatFormConnected")

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
	console.log("appendOwnMessage");

	const listItem = document.createElement("li");

	listItem.className = "message message-own";
	listItem.dataset.messageId = message.message_id;

	const body = document.createElement("div");
	body.className = "message-body";

	const text = document.createElement("span");
	text.className = "message-text";
    appendTextWithLinks(text, message.body);

	const footer = document.createElement("div");
	footer.className = "message-body-footer";
	footer.textContent = message.display_time;

	body.appendChild(text);
	body.appendChild(footer);
	listItem.appendChild(body);

	configureSingleMessagePressHandler(listItem);

	messagesElement.appendChild(listItem);
}

function appendOtherMessage(messagesElement, message) {
	console.log("appendOtherMessage");

	const listItem = document.createElement("li");

	listItem.className = "message test-messsage-other";
	listItem.dataset.messageId = message.message_id;

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

	const text = document.createElement("span");
	text.className = "message-text";
    appendTextWithLinks(text, message.body);

	const footer = document.createElement("div");
	footer.className = "message-body-footer";
	footer.textContent = message.display_time;

	body.appendChild(text);
	body.appendChild(footer);

	content.appendChild(header);
	content.appendChild(body);

	listItem.appendChild(pictureBox);
	listItem.appendChild(content);

	configureSingleMessagePressHandler(listItem);

	messagesElement.appendChild(listItem);
}

function handleIncomingMessage(event) {
	console.log("handleIncomingMessage")
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

	const currentUserId = messagesElement.dataset.currentUserId;

	if (message.user_id === currentUserId) {
		appendOwnMessage(messagesElement, message);
	} else {
		appendOtherMessage(messagesElement, message);
	}

	messagesElement.scrollTop = messagesElement.scrollHeight;
}

function connectChatWebSocket() {
	console.log("connectChatWebSocket")
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
            reconnectTimer = null;
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
		const chatID = form.dataset.chatId;

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
		console.log("chat-id: ", chatID)

		window.chatSocket.send(
			JSON.stringify({
				type: "chat_message",
				chat_id: chatID,
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

function addPressHandlers(element, {
    holdDuration = 600,
    onTap,
    onHold,
}) {
    let holdTimer = null;
    let holdTriggered = false;
    let startX = 0;
    let startY = 0;

    function cancelHoldTimer() {
        if (holdTimer) {
            clearTimeout(holdTimer);
            holdTimer = null;
        }
    }

	function cancelPress() {
		cancelHoldTimer();
		holdTriggered = false;
	}

    element.addEventListener("pointerdown", (event) => {
		if (event.target.closest("a")) {
			return;
		}

        if (event.pointerType === "mouse" && event.button !== 0) {
            return;
        }

        holdTriggered = false;
        startX = event.clientX;
        startY = event.clientY;

        holdTimer = setTimeout(() => {
            holdTriggered = true;

            navigator.vibrate?.(40);

            onHold?.(event, element);
        }, holdDuration);
    });

    element.addEventListener("pointermove", (event) => {
        const distanceX = Math.abs(event.clientX - startX);
        const distanceY = Math.abs(event.clientY - startY);

        if (distanceX > 10 || distanceY > 10) {
            cancelHoldTimer();
        }
    });

    element.addEventListener("pointerup", (event) => {
        cancelHoldTimer();

		if (event.target.closest("a")) {
			holdTriggered = false;
			return;
		}

		if (!holdTriggered) {
            onTap?.(event, element);
        }

        holdTriggered = false;
    });

    element.addEventListener("pointercancel", cancelPress);
    element.addEventListener("pointerleave", cancelPress);

    element.addEventListener("contextmenu", (event) => {
		if (event.target.closest("a")) {
			return;
		}
        event.preventDefault();
    });
}






const chatBackdrop = document.getElementById("chat-backdrop");
const chatPopupMessage = document.getElementById("chat-popup-message")

function toggleChatPopupMenu(event) {
	console.log("toggleChatPopupMenu")
	event.preventDefault();
	event.stopPropagation();

	const options = event.currentTarget.closest(".chat-reactions-box");

	document.querySelectorAll(".chat-reactions-box.open").forEach(menu => {
		console.log("querySelectorAll - .chat-reactions-box.open")
		if (menu !== options) {
			menu.classList.remove("open");
			chatBackdrop?.classList.remove("open");
			chatPopupMessage?.classList.remove("open");
		}
	});
}

function closeChatReactionPopup(event) {
	console.log("closeItemCardMenu")
	const chatBackdrop = document.getElementById("chat-backdrop");
	const chatPopupMessage = document.getElementById("chat-popup-message")
	document
		.querySelectorAll(".chat-reactions-box.open")
		.forEach((menu) => {
			menu.classList.remove("open");
			chatPopupMessage?.classList.remove("open");
			chatBackdrop?.classList.remove("open");
		});
		chatBackdrop?.classList.remove("open")
	}

function openMessageOptions(messageId, messageElement) {
	console.log("openMessageOptions")
	// const chatBackdrop = document.getElementById("chat-backdrop");
	const chatReactionPopup = document.getElementById("chat-reactions-box");
	const chatPopupMessage = document.getElementById("chat-popup-message");

	console.log("Show options for message:", messageId);

	document
		.querySelectorAll(".message-selected")
		.forEach((element) => {
			element.classList.remove("message-selected");
		});

	chatBackdrop?.classList.add("open");
	chatReactionPopup?.classList.add("open");
	chatPopupMessage?.classList.add("open");
	messageElement.classList.add("message-selected");

	chatBackdrop?.addEventListener("click", () => {
		console.log("event listener: chatBackdrop")
        closeChatReactionPopup();
		chatBackdrop?.classList.remove("open");
      });

}

function closeMessageOptions() {
	console.log("closeMessageOptions")
	document
		.querySelectorAll(".message-selected")
		.forEach((element) => {
			element.classList.remove("message-selected");
			chatBackdrop?.classList.remove("open");
		});
}







function configureMessagePressHandlers() {
	const messages = document.querySelectorAll(".message");

	for (const messageElement of messages) {
		configureSingleMessagePressHandler(messageElement)
	}
}

function configureSingleMessagePressHandler(messageElement) {
	console.log("configureSingleMessagePressHandler")
	addPressHandlers(messageElement, {
		holdDuration: 600,

		onTap: (_event, element) => {
			console.log(
				"Tapped message",
				element.dataset.messageId,
			);
		},

		onHold: (_event, element) => {
			openMessageOptions(
				element.dataset.messageId,
				element,
			);
		},
	});
}



configureMessageForm();
configureMessagePressHandlers();
linkifyExistingMessages();
scrollMessagesToBottom();




function closeAddMenu() {
	addMenu?.classList.remove("open");
	sideMenuBackdrop?.classList.remove("open");
}
