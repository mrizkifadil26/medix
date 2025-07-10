let allIcons = [];

async function loadUnusedIcons() {
  try {
    const res = await fetch("data/reports/unused-icons.report.json");
    const json = await res.json();
    renderUnusedIcons(json.groups);
  } catch (err) {
    console.error("❌ Failed to load unused icon report:", err);
  }
}

function renderUnusedIcons(groups) {
  const container = document.getElementById("unused-icons-container");
  if (!groups || Object.keys(groups).length === 0) {
    container.innerHTML = "<p>✅ All icons are in use.</p>";
    return;
  }

  let html = "";
  const sortedGroups = Object.keys(groups).sort();

  for (const group of sortedGroups) {
    const icons = groups[group];
    html += `
      <div class="icon-group">
        <div class="group-header" onclick="toggleGroup(this)">
          <span class="group-name">${GENRE_EMOJIS[group]} ${group}</span>
          <span class="group-count">${icons.length} icons</span>
        </div>
        <ul class="icon-list" style="display: none">
          ${icons.map(icon => `
            <li class="icon-row">
              <span class="icon-name">${removeExtension(icon.name)}</span>
              <span class="icon-source tag ${icon.source}">${icon.source}</span>
            </li>
          `).join("")}
        </ul>
      </div>
    `;
  }

  container.innerHTML = html;
}

function removeExtension(name) {
  return name.replace(/\.ico$/i, "");
}

function toggleGroup(header) {
  const list = header.nextElementSibling;
  if (list) {
    list.style.display = list.style.display === "none" ? "block" : "none";
  }
}

// document.getElementById("icon-search").addEventListener("input", e => {
//   const query = e.target.value.toLowerCase();
//   document.querySelectorAll("details").forEach(d => {
//     const group = d.querySelector(".unused-icons-group");
//     if (!d.open || !group) return;

//     const rows = group.querySelectorAll("tbody tr");
//     rows.forEach(row => {
//       const text = row.textContent.toLowerCase();
//       row.style.display = text.includes(query) ? "" : "none";
//     });
//   });
// });

window.addEventListener("DOMContentLoaded", loadUnusedIcons);
