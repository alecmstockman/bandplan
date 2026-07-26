
console.log("NEW chat-websocket.js loaded");

function connectChatWebSocket() {
    console.log("- connectChatWebSocket")
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const socketURL = `${protocol}//${window.location.host}/ws/chat`;

    console.log("WebSocket URL:", socketURL);
    const socket = new WebSocket(socketURL);

    socket.addEventListener("open", () => {
        console.log("   Chat WebSocket connected");
    });

    socket.addEventListener("message", (event) => {
        const message = JSON.parse(event.data)
		console.log("WebSocket message received:", event.data);
	});

	socket.addEventListener("close", (event) => {
		console.log("Chat WebSocket disconnected:", event.code);
	});

	socket.addEventListener("error", (error) => {
		console.error("Chat WebSocket error:", error);
	});

	return socket;
}

if (
	!window.chatSocket ||
	window.chatSocket.readyState === WebSocket.CLOSED
) {
	window.chatSocket = connectChatWebSocket();
}


