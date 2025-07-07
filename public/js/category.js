function getCategoryType() {
  const filename = window.location.pathname.split("/").pop();
  if (filename.includes("movies")) return "movies";
  if (filename.includes("tvshows")) return "tv_shows";
  return null;
}

async function loadCategory() {
  const type = getCategoryType();
  if (!type) {
    document.getElementById("category-title").innerText = "Unknown Category";
    return;
  }

  const response = await fetch(`data/${type}.json`);
  const genres = await response.json();

  document.getElementById("category-title").innerText =
    type === "movies" ? "🎬 Movies by Genre" : "📺 TV Shows by Genre";

  const container = document.getElementById("category-container");
  container.innerHTML = "";

  genres.forEach(genre => {
    const section = document.createElement("section");
    const heading = document.createElement("h2");
    heading.innerText = genre.genre;
    section.appendChild(heading);

    const ul = document.createElement("ul");
    genre.titles.forEach(title => {
      const li = document.createElement("li");
      li.textContent = title.name;

      const statusIcon = document.createElement("span");
      statusIcon.style.marginLeft = "5px";
      statusIcon.style.cursor = "help";

      switch (title.status) {
        case "warn":
          statusIcon.textContent = "⚠️";
          statusIcon.title = "Missing Thumbnail";
          break;
        case "missing":
          statusIcon.textContent = "❌";
          statusIcon.title = "Missing Icon";
          break;
        default:
          statusIcon.textContent = "✅";
          statusIcon.title = "Complete";
          break;
      }

      li.appendChild(statusIcon);
      ul.appendChild(li);
    });

    section.appendChild(ul);
    container.appendChild(section);
  });
}

window.addEventListener("DOMContentLoaded", loadCategory);
