async function loadProgress() {
  const res = await fetch("data/progress.json");
  const json = await res.json();

  renderProgressBar(json.percent, json.done, json.total);
  renderProgressTable(json.genres);
  renderProgressCards(json.genres);
}

function renderProgressBar(percent, done, total) {
  const container = document.getElementById("progress-bar-container");
  container.innerHTML = `
    <div class="progress-container">
      <strong>Progress:</strong>
      <div class="progress-bar">
        <div class="progress-fill" style="width: ${percent}%;"></div>
        <div class="progress-text">${percent}% (${done}/${total})</div>
      </div>
    </div>
  `;
}

function renderProgressTable(genres) {
  const container = document.getElementById("progress-table");

  let html = `
    <table class="progress-table">
      <thead>
        <tr>
          <th>Genre</th>
          <th>Progress</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
  `;

  genres.forEach(g => {
    const raw = g.raw || 1;
    const ico = g.ico || 0;
    const png = g.png || 0;
    const pngOnly = png - ico;
    const rawOnly = raw - png;

    const icoWidth = (ico / raw) * 100;
    const pngWidth = (pngOnly / raw) * 100;
    const rawWidth = (rawOnly / raw) * 100;

    const percent = Math.round((ico / raw) * 100);

    html += `
      <tr>
        <td data-label="Genre">
          <span class="genre-label">
            ${g.icon || "🎬"} ${g.genre}
          </span>
        </td>

        <td data-label="Progress">
          <div class="bar-wrapper">
            <div class="bar-container">
              <div class="bar-ico" style="width: ${icoWidth}%"></div>
              <div class="bar-png" style="width: ${pngWidth}%"></div>
              <div class="bar-raw" style="width: ${rawWidth}%"></div>
              <div class="bar-label">${percent}%</div>
            </div>
            <div class="bar-info">
              <span class="bar-count ico">${ico}</span> /
              <span class="bar-count png">${png}</span> /
              <span class="bar-count raw">${raw}</span>
            </div>
          </div>
        </td>
        <td data-label="Status">${g.status}</td>
      </tr>
    `;
  });

  html += `
      </tbody>
    </table>
    <div class="legend">
      <strong>Legend:</strong>
      <span class="box bar-ico"></span> ICO
      <span class="box bar-png"></span> PNG
      <span class="box bar-raw"></span> RAW
    </div>
  `;

  container.innerHTML = html;
}

function renderProgressCards(genres) {
  const container = document.getElementById("progress-cards");
  container.innerHTML = "";

  genres.forEach(g => {
    const raw = g.raw || 1;
    const ico = g.ico || 0;
    const png = g.png || 0;
    const pngOnly = png - ico;
    const rawOnly = raw - png;

    const icoWidth = (ico / raw) * 100;
    const pngWidth = (pngOnly / raw) * 100;
    const rawWidth = (rawOnly / raw) * 100;

    const percent = Math.round((ico / raw) * 100);

    const card = document.createElement("div");
    card.className = "genre-card";
    card.innerHTML = `
      <div class="genre-header">
        <div class="genre-name">${g.icon || "🎬"} ${g.genre}</div>
        <div class="genre-status">${g.status}</div>
      </div>
      <div class="bar-container">
        <div class="bar-ico" style="width: ${icoWidth}%;"></div>
        <div class="bar-png" style="width: ${pngWidth}%;"></div>
        <div class="bar-raw" style="width: ${rawWidth}%;"></div>
        <div class="bar-label">${Math.round(percent)}%</div>
      </div>

      <div class="bar-info">
        <span class="bar-count ico">${ico}</span> /
        <span class="bar-count png">${png}</span> /
        <span class="bar-count raw">${raw}</span>
      </div>
    `;
    container.appendChild(card);
  });
}


window.addEventListener("DOMContentLoaded", loadProgress);
