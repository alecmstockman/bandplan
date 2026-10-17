

const toastMessages = {
    "song-saved": "Song saved successfully",
    "song-deleted": "Song deleted successfully",
    "event-created": "Event created successfully",
};

const url = new URL(window.location.href);
const toastCode = url.searchParams.get("toast");

if (toastCode !== null) {
    if (Object.hasOwn(toastMessages, toastCode)) {
        showToast(toastMessages[toastCode]);
    }

    url.searchParams.delete("toast");
    window.history.replaceState(window.history.state, "", url);
}


function showToast(message) {
    console.log("showToast")
    const container = document.getElementById("toast-container");
    if (!container) {
        console.error("Toast container was not found");
        return;
    }

    const toast = document.createElement("div");
    toast.classList.add("toast-card");
    toast.textContent = message;

    container.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 4000)
}

document.addEventListener("showToast", (event) => {
    console.log("showToast Event Listener")
    showToast(event.detail.message);
});