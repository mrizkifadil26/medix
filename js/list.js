const GENRE_EMOJIS = {
  "Horror": "👻",
  "Drama": "🎭",
  "Animation": "🎨",
  "Comedy": "😂",
  "Thriller": "🔪",
  "Action": "🔥",
  "Romance": "💕",
  "Sci-Fi": "👽",
  "Fantasy": "🧙",
  "Family": "🏠",
  "Mystery": "🧩",
  "Crime": "🔍"
};

async function loadGroupedTitles() {
  const res = await fetch("data/movies.json");
  const data = await res.json();
  renderGroupedList(data);
}

function renderGroupedList(genres) {
  const container = document.getElementById("genre-list");
  container.innerHTML = "";

  Object.entries(genres).forEach(([genre, items]) => {
    const genreDetails = document.createElement("details");
    genreDetails.className = "collapsible";

    const genreSummary = document.createElement("summary");
    const icon = GENRE_EMOJIS[genre] || "🎬";

    genreSummary.innerHTML = `
      <strong>${icon} ${genre}</strong>
      <span class="count">(${items.length} items)</span>
      <span class="genre-chevron chevron" style="float:right;">${getChevronDown()}</span>
    `;
    genreDetails.appendChild(genreSummary);

    genreDetails.addEventListener("toggle", () => {
      const chevron = genreDetails.querySelector(".genre-chevron");
      chevron.innerHTML = genreDetails.open ? getChevronUp() : getChevronDown();
    });

    const listWrapper = document.createElement("div");
    listWrapper.className = "collapsible-body";

    const ul = document.createElement("ul");
    ul.className = "title-list";

    items.forEach(item => {
      const li = document.createElement("li");

      if (item.group && Array.isArray(item.group)) {
        const collectionDetails = document.createElement("details");

        const collectionSummary = document.createElement("summary");
        collectionSummary.innerHTML = `
          📂 ${item.name}
          <span class="group-chevron chevron" style="float:right;">${getChevronDown()}</span>
        `;
        collectionDetails.appendChild(collectionSummary);

        collectionDetails.addEventListener("toggle", () => {
          const chevron = collectionDetails.querySelector(".group-chevron");
          chevron.innerHTML = collectionDetails.open ? getChevronUp() : getChevronDown();
        });

        const subList = document.createElement("ul");
        subList.className = "title-list";

        item.group.forEach(sub => {
          const subLi = document.createElement("li");
          subLi.innerHTML = `${getStatusEmoji(sub.status)} ${sub.name}`;
          subList.appendChild(subLi);
        });

        collectionDetails.appendChild(subList);
        li.appendChild(collectionDetails);
      } else {
        li.innerHTML = `${getStatusEmoji(item.status)} ${item.name}`;
      }

      ul.appendChild(li);
    });

    listWrapper.appendChild(ul);
    genreDetails.appendChild(listWrapper);
    container.appendChild(genreDetails);
  });
}

function getStatusEmoji(status) {
  switch (status) {
    case "ok": return "✅";
    case "warn": return "⚠️";
    case "missing": return "❌";
    default: return "";
  }
}

// SVG Chevron helpers
function getChevronDown() {
  return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
    stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>`;
}

function getChevronUp() {
  return `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
    stroke-linecap="round" stroke-linejoin="round"><polyline points="18 15 12 9 6 15"/></svg>`;
}

window.addEventListener("DOMContentLoaded", loadGroupedTitles);
