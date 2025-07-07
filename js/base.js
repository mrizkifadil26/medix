document.addEventListener("DOMContentLoaded", function () {
    // 🔥 Highlight active navbar link
    const currentFile = window.location.pathname.split("/").pop() || "index.html";
    document.querySelectorAll("header nav a").forEach(link => {
        if (link.getAttribute("href") === currentFile) {
            link.classList.add("active");
        }
    });
});