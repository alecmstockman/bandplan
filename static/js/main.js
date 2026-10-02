

console.log("main.js loaded");


function getCSRFToken() {
    const tokenCookie = document.cookie
        .split("; ")
        .find((cookie) => cookie.startsWith("csrf_token="));

    return tokenCookie
        ? decodeURIComponent(tokenCookie.slice("csrf_token=".length))
        : "";
}


document.addEventListener("htmx:configRequest", (event) => {
    const unsafeMethods = ["post", "put", "patch", "delete"];

    if (!unsafeMethods.includes(event.detail.verb.toLowerCase())) {
        return;
    }

    const token = getCSRFToken();

    if (token) {
        event.detail.headers["X-CSRF-Token"] = token;
    }
});



async function copyToClipboard(text, button) {
    try {
        await navigator.clipboard.writeText(text);

        button.classList.add("copied");

        setTimeout(() => {
            button.classList.remove("copied");
        }, 1200);

    } catch (err) {
        console.error(err);
    }
}


document.addEventListener("click", (event) => {
    const target = event.target;

    if (!(target instanceof Element)) {
        return;
    }

    const leftMenu = document.getElementById("left-menu");
    const rightMenu = document.getElementById("right-menu");

    if (target.closest("#menu-button")) {
        leftMenu?.classList.toggle("toggle");
        rightMenu?.classList.remove("toggle");
        return;
    }

    if (target.closest("#close-left-menu")) {
        leftMenu?.classList.remove("toggle");
        return;
    }

    if (target.closest("#right-menu-button")) {
        rightMenu?.classList.toggle("toggle");
        leftMenu?.classList.remove("toggle");
        return;
    }

    if (target.closest("#close-right-menu")) {
        rightMenu?.classList.remove("toggle");
        return;
    }

    if (
        leftMenu?.classList.contains("toggle") &&
        !target.closest("#left-menu")
    ) {
        leftMenu.classList.remove("toggle");
    }

    if (
        rightMenu?.classList.contains("toggle") &&
        !target.closest("#right-menu")
    ) {
        rightMenu.classList.remove("toggle");
    }
});

document.addEventListener("keydown", (event) => {
    if (event.key !== "Escape") {
        return;
    }

    document.getElementById("left-menu")
        ?.classList.remove("toggle");

    document.getElementById("right-menu")
        ?.classList.remove("toggle");
});
