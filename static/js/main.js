


// async function copyToClipboard(text, button) {
//     try {
//         await navigator.clipboard.writeText(text);

//         button.classList.add("copied");

//         setTimeout(() => {
//             button.classList.remove("copied");
//         }, 1200);

//     } catch (err) {
//         console.error("Copy failed:", err);
//     }
// }




console.log("main.js loaded");

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