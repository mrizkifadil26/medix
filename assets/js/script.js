document.querySelectorAll(".toggle-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    const section = btn.closest(".dashboard-section");
    const content = section.querySelector(".dashboard-content");

    const isCollapsed = content.classList.toggle("collapsed");
    const expanded = !isCollapsed;

    btn.setAttribute("aria-expanded", expanded);
    btn.innerHTML = expanded ? "🔽 Collapse" : "▶️ Show All";
  });
});
